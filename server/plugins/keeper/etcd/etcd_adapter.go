package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"ivory/clients/etcd"
	"ivory/plugins/keeper"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var ErrScheduleNotSupported = errors.New("scheduled switchover is not supported by etcd")
var ErrCandidateRequired = errors.New("candidate is required, etcd cannot choose a leader randomly")
var ErrCandidateNotFound = errors.New("candidate was not found among cluster members")

// maxApplyGap is etcd's own maxGapBetweenApplyAndCommitIndex: past it the
// member refuses new requests with "too many requests".
const maxApplyGap = 5000

// NOTE: validate that is matches interface in compile-time
var _ keeper.Adapter = (*Plugin)(nil)

// Plugin manages an etcd cluster as the keeper entity. The keeper
// connection host/port is the etcd client endpoint (keeperPort == dbPort ==
// client port convention). Only overview and leader move are supported;
// everything else is excluded from SupportedFeatures.
type Plugin struct{}

func NewPlugin() *Plugin {
	return &Plugin{}
}

// List asks the node Ivory reached about itself and the leader only. Peers are
// reported from the member list alone, so each one's own List supplies its
// state: probing every member from every node fetched the same answers N times.
func (p *Plugin) List(request keeper.Request) ([]keeper.Response, int, error) {
	ctx, client, err := p.connect(request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer client.Close()

	var memberList *clientv3.MemberListResponse
	var selfStatus *clientv3.StatusResponse
	var errMembers, errStatus error
	var wg sync.WaitGroup
	wg.Go(func() { memberList, errMembers = client.MemberList(ctx) })
	wg.Go(func() { selfStatus, errStatus = client.StatusSelf(ctx) })
	wg.Wait()
	if err := errors.Join(errMembers, errStatus); err != nil {
		return nil, http.StatusBadRequest, err
	}

	members := make([]member, 0, len(memberList.Members))
	for _, m := range memberList.Members {
		members = append(members, member{ID: m.ID, Name: m.Name, ClientURLs: m.ClientURLs, IsLearner: m.IsLearner})
	}

	selfID := selfStatus.Header.GetMemberId()
	self := mapStatus(selfStatus)
	leader, leaderURL := findLeader(members, selfID, self)
	if leaderURL != "" {
		leader = p.probeLeader(client, leader, leaderURL)
	}

	responses := mapMembers(members, selfID, self, leader)
	// NOTE: the cluster id belongs to the reply as a whole rather than to any
	// one member, which is why it is applied here instead of inside mapMember.
	// Enumerating members already catches a node borrowed from another etcd of
	// several - its own peers show up as nodes nobody configured - but not one
	// borrowed from a single-node cluster, which reports itself and nothing
	// else. The id is what contradicts that one.
	if clusterId := memberList.Header.GetClusterId(); clusterId != 0 {
		identity := strconv.FormatUint(clusterId, 10)
		for i := range responses {
			responses[i].DiscoveredCluster = &identity
		}
	}
	return responses, http.StatusOK, nil
}

// probeLeader reads the leader's raft index, which a follower's own status
// lacks and its lag is measured against. It has its own timeout so a slow
// leader cannot spend the budget the asked node's own reply needed.
func (p *Plugin) probeLeader(client *connection, leader leaderStatus, url string) leaderStatus {
	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()
	status, err := client.Status(ctx, url)
	if err != nil {
		leader.Err = fmt.Errorf("leader is unreachable at %s: %w", url, err)
		return leader
	}
	leader.RaftIndex = status.RaftIndex
	return leader
}

func (p *Plugin) Switchover(request keeper.Request) (*string, int, error) {
	body, errBody := parseSwitchoverBody(request.Body)
	if errBody != nil {
		return nil, http.StatusBadRequest, errBody
	}

	ctx, client, err := p.connect(request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer client.Close()

	memberList, errMembers := client.MemberList(ctx)
	if errMembers != nil {
		return nil, http.StatusBadRequest, errMembers
	}
	members := make([]member, 0, len(memberList.Members))
	for _, m := range memberList.Members {
		members = append(members, member{ID: m.ID, Name: m.Name, ClientURLs: m.ClientURLs, IsLearner: m.IsLearner})
	}

	candidateID, errCandidate := resolveCandidate(members, *body.Candidate)
	if errCandidate != nil {
		return nil, http.StatusBadRequest, errCandidate
	}

	if _, errMove := client.MoveLeader(ctx, candidateID); errMove != nil {
		return nil, http.StatusBadRequest, errMove
	}
	response := "leader moved to " + *body.Candidate
	return &response, http.StatusOK, nil
}

func (p *Plugin) Config(keeper.Request) (any, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) ConfigUpdate(keeper.Request) (any, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) DeleteSwitchover(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) Reinitialize(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) Restart(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) DeleteRestart(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) Reload(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) Failover(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) Activate(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (p *Plugin) Pause(keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

// connection bundles the etcd client with the context bound to its request
// timeout, so callers get both from connect and clean up both with one
// deferred Close.
type connection struct {
	*etcd.Client
	cancel context.CancelFunc
}

func (c *connection) Close() {
	c.cancel()
	err := c.Client.Close()
	if err != nil {
		slog.Error("failed to close etcd client", "error", err)
	}
}

// connect opens an etcd client for the request's endpoint and returns a
// context bound to the client's request timeout together with it.
func (p *Plugin) connect(request keeper.Request) (context.Context, *connection, error) {
	var username, password string
	if request.Credentials != nil {
		username = request.Credentials.Username
		password = request.Credentials.Password
	}
	client, err := etcd.Connect(etcd.Config{
		Endpoints: []string{request.Host + ":" + strconv.Itoa(request.Port)},
		Username:  username,
		Password:  password,
		TLS:       request.TlsConfig,
	})
	if err != nil {
		return nil, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	return ctx, &connection{Client: client, cancel: cancel}, nil
}

func parseSwitchoverBody(body any) (*switchoverBody, error) {
	bytes, errMarshal := json.Marshal(body)
	if errMarshal != nil {
		return nil, errMarshal
	}
	var parsed switchoverBody
	if errUnmarshal := json.Unmarshal(bytes, &parsed); errUnmarshal != nil {
		return nil, errUnmarshal
	}
	if parsed.ScheduledAt != nil && *parsed.ScheduledAt != "" {
		return nil, ErrScheduleNotSupported
	}
	if parsed.Candidate == nil || *parsed.Candidate == "" {
		return nil, ErrCandidateRequired
	}
	return &parsed, nil
}

func resolveCandidate(members []member, candidate string) (uint64, error) {
	for _, m := range members {
		if m.Name == candidate {
			return m.ID, nil
		}
		for _, clientUrl := range m.ClientURLs {
			parsed, err := url.Parse(clientUrl)
			if err != nil {
				continue
			}
			if parsed.Hostname() == candidate {
				return m.ID, nil
			}
		}
	}
	return 0, fmt.Errorf("%w: %s", ErrCandidateNotFound, candidate)
}

func mapStatus(status *clientv3.StatusResponse) endpointStatus {
	return endpointStatus{
		Leader:           status.Leader,
		RaftIndex:        status.RaftIndex,
		RaftAppliedIndex: status.RaftAppliedIndex,
		RaftTerm:         status.RaftTerm,
		Version:          status.Version,
		Errors:           status.Errors,
	}
}

// findLeader takes the leader from the asked member's own status. It returns
// the url to probe for the leader's raft index, or none when the asked member
// is the leader itself, sees no leader, or cannot address the one it sees.
func findLeader(members []member, selfID uint64, self endpointStatus) (leaderStatus, string) {
	leader := leaderStatus{ID: self.Leader}
	if self.Leader == 0 {
		return leader, ""
	}
	if self.Leader == selfID {
		leader.RaftIndex = self.RaftIndex
		return leader, ""
	}
	for _, m := range members {
		if m.ID == self.Leader && len(m.ClientURLs) > 0 {
			return leader, m.ClientURLs[0]
		}
	}
	return leader, ""
}

func mapMembers(members []member, selfID uint64, self endpointStatus, leader leaderStatus) []keeper.Response {
	responses := make([]keeper.Response, 0, len(members))
	for _, m := range members {
		if m.ID != selfID {
			responses = append(responses, mapPeer(m, leader.ID))
			continue
		}
		response := mapMember(m, self, leader.ID, leader.RaftIndex)
		if leader.Err != nil {
			response.Warnings = append(response.Warnings, leader.Err.Error())
		}
		responses = append(responses, response)
	}
	return responses
}

// mapPeer reports a member nobody contacted: its role and address are the
// asked member's topology, while its state and lag are left for its own List.
func mapPeer(m member, leaderID uint64) keeper.Response {
	var role keeper.Role = keeper.Unknown
	if leaderID != 0 {
		role = keeper.Replica
		if m.ID == leaderID {
			role = keeper.Leader
		}
	}
	state := keeper.StateUnknown
	if len(m.ClientURLs) == 0 {
		state = keeper.StateStarting
	}
	response := mapMember(m, endpointStatus{}, 0, 0)
	response.State = state
	response.Role = role
	response.Lag = -1
	return response
}

// mapTags reports the per-member facts etcd's own status carries that
// Role/State/Lag cannot express: a learner does not vote and can never be
// elected, a version that differs from its peers' is a half-finished rolling
// upgrade, and a raft term behind the rest is a member that missed an election.
// All three come off the Status reply that Role and State are already read
// from, so none of them costs a request. The db size is not here: Ivory does
// not know the backend quota, so the number alone says nothing about the one
// failure it would warn of. An unreachable member reports a zero status, and
// every field is omitted rather than printed as a zero.
func mapTags(m member, status endpointStatus) *map[string]any {
	tags := map[string]any{}
	if m.IsLearner {
		tags["learner"] = true
	}
	if status.Version != "" {
		tags["version"] = status.Version
	}
	if status.RaftTerm > 0 {
		tags["raftTerm"] = status.RaftTerm
	}
	if len(tags) == 0 {
		return nil
	}
	return &tags
}

// mapWarnings reports what the member says about itself and why Ivory could not
// ask it. The url is named because it is the one the member advertised, so an
// address Ivory cannot route to reads differently from a member that is down.
func mapWarnings(m member, status endpointStatus) []string {
	if status.Err != nil {
		if len(m.ClientURLs) > 0 {
			return []string{"member is unreachable at " + m.ClientURLs[0] + ": " + status.Err.Error()}
		}
		return []string{"member is unreachable: " + status.Err.Error()}
	}
	var warnings []string
	for _, e := range status.Errors {
		warnings = append(warnings, "member reports: "+e)
	}
	if status.RaftIndex > status.RaftAppliedIndex && status.RaftIndex-status.RaftAppliedIndex > maxApplyGap {
		warnings = append(warnings, fmt.Sprintf(
			"%d committed entries are not applied yet, past %d etcd refuses new requests",
			status.RaftIndex-status.RaftAppliedIndex, maxApplyGap))
	}
	return warnings
}

func mapMember(m member, status endpointStatus, leaderID uint64, leaderRaftIndex uint64) keeper.Response {
	name := m.Name
	var keeperStatus keeper.Status = keeper.Active

	state := keeper.StateRunning
	var role keeper.Role = keeper.Unknown
	// NOTE: lag is the raft index difference to the leader, unlike patroni's lag value
	lag := int64(-1)
	switch {
	case len(m.ClientURLs) == 0:
		state = keeper.StateStarting
	case status.Err != nil:
		state = keeper.StateUnreachable
	case m.ID == leaderID:
		role = keeper.Leader
		lag = 0
	default:
		role = keeper.Replica
		if leaderRaftIndex >= status.RaftIndex {
			lag = int64(leaderRaftIndex - status.RaftIndex)
		}
	}

	response := keeper.Response{
		Key:            &name,
		Status:         &keeperStatus,
		State:          state,
		Role:           role,
		Lag:            lag,
		Tags:           mapTags(m, status),
		DiscoveredName: &name,
	}
	response.Warnings = mapWarnings(m, status)

	if len(m.ClientURLs) > 0 {
		if parsed, err := url.Parse(m.ClientURLs[0]); err == nil {
			host := parsed.Hostname()
			if port, errPort := strconv.Atoi(parsed.Port()); errPort == nil {
				response.DiscoveredHost = &host
				response.DiscoveredKeeperPort = &port
				response.DiscoveredDbPort = &port
			}
		}
	}

	return response
}
