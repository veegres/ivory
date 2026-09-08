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
)

var ErrScheduleNotSupported = errors.New("scheduled switchover is not supported by etcd")
var ErrCandidateRequired = errors.New("candidate is required, etcd cannot choose a leader randomly")
var ErrCandidateNotFound = errors.New("candidate was not found among cluster members")

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

func (p *Plugin) List(request keeper.Request) ([]keeper.Response, int, error) {
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
	statuses := make(map[uint64]endpointStatus, len(memberList.Members))
	// NOTE: a single unreachable member must not fail the whole overview, but
	// its error must still make it back to the caller (not just the
	// "unreachable" state) so it ends up in the node's warnings.
	var errs error
	for _, m := range memberList.Members {
		members = append(members, member{ID: m.ID, Name: m.Name, ClientURLs: m.ClientURLs, IsLearner: m.IsLearner})
		if len(m.ClientURLs) == 0 {
			continue
		}
		status, errStatus := client.Status(ctx, m.ClientURLs[0])
		if errStatus != nil {
			statuses[m.ID] = endpointStatus{Err: errStatus}
			errs = errors.Join(errs, fmt.Errorf("member %q is unreachable: %w", m.Name, errStatus))
			continue
		}
		statuses[m.ID] = endpointStatus{
			Leader:    status.Leader,
			RaftIndex: status.RaftIndex,
			RaftTerm:  status.RaftTerm,
			Version:   status.Version,
			DbSize:    status.DbSize,
		}
	}

	status := http.StatusOK
	if errs != nil {
		status = http.StatusServiceUnavailable
	}

	responses := mapMembers(members, statuses)
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
	return responses, status, errs
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
		slog.Error("failed to close ssh client", "error", err)
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

func mapMembers(members []member, statuses map[uint64]endpointStatus) []keeper.Response {
	leaderID, leaderRaftIndex := findLeader(statuses)

	responses := make([]keeper.Response, 0, len(members))
	for _, m := range members {
		responses = append(responses, mapMember(m, statuses[m.ID], leaderID, leaderRaftIndex))
	}
	return responses
}

// findLeader only designates a leader when a strict majority of the
// responding members agree on the same leader ID. During an election or a
// network partition, reachable members can briefly disagree (some still
// reporting a stale leader, others none yet); without this check, whichever
// member's status happened to be read first (map iteration order is
// randomized) would arbitrarily decide the reported leader, causing it to
// flicker between candidates across polls even though nothing in the
// cluster actually changed. Requiring majority agreement makes the result
// deterministic: either the members agree and the winner is unambiguous, or
// they don't and no leader is reported (surfacing the cluster's "no leader
// found" warning instead of a guess).
func findLeader(statuses map[uint64]endpointStatus) (uint64, uint64) {
	votes := make(map[uint64]int)
	for _, status := range statuses {
		if status.Err == nil && status.Leader != 0 {
			votes[status.Leader]++
		}
	}
	var leaderID uint64
	var leaderVotes int
	for id, count := range votes {
		if count > leaderVotes {
			leaderID = id
			leaderVotes = count
		}
	}
	if leaderVotes*2 <= len(statuses) {
		return 0, 0
	}
	if leader, ok := statuses[leaderID]; ok && leader.Err == nil {
		return leaderID, leader.RaftIndex
	}
	return leaderID, 0
}

// mapTags reports the per-member facts etcd's own status carries that
// Role/State/Lag cannot express: a learner does not vote and can never be
// elected, a version that differs from its peers' is a half-finished rolling
// upgrade, a db size approaching the backend quota is what makes a cluster go
// read-only, and a raft term behind the rest is a member that missed an
// election. An unreachable member reports a zero status, and every field is
// omitted rather than printed as a zero.
func mapTags(m member, status endpointStatus) *map[string]any {
	tags := map[string]any{}
	if m.IsLearner {
		tags["learner"] = true
	}
	if status.Version != "" {
		tags["version"] = status.Version
	}
	if status.DbSize > 0 {
		tags["dbSize"] = formatDbSize(status.DbSize)
	}
	if status.RaftTerm > 0 {
		tags["raftTerm"] = status.RaftTerm
	}
	if len(tags) == 0 {
		return nil
	}
	return &tags
}

// formatDbSize reports KiB below a megabyte: a freshly bootstrapped etcd holds
// about 20 KiB, which as megabytes reads "0.0 MiB" and looks like a node that
// failed to answer rather than a healthy empty one.
func formatDbSize(bytes int64) string {
	if bytes >= 1024*1024 {
		return fmt.Sprintf("%.1f MiB", float64(bytes)/(1024*1024))
	}
	return fmt.Sprintf("%.0f KiB", float64(bytes)/1024)
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
