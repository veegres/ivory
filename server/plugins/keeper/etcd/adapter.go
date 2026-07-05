package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"ivory/clients/etcd"
	"ivory/plugins/keeper"
	"net/http"
	"net/url"
	"strconv"
)

var ErrScheduleNotSupported = errors.New("scheduled switchover is not supported by etcd")
var ErrCandidateRequired = errors.New("candidate is required, etcd cannot choose a leader randomly")
var ErrCandidateNotFound = errors.New("candidate was not found among cluster members")

// NOTE: validate that is matches interface in compile-time
var _ keeper.Adapter = (*Adapter)(nil)

// Adapter manages an etcd cluster as the keeper entity. The keeper
// connection host/port is the etcd client endpoint (keeperPort == dbPort ==
// client port convention). Only overview and leader move are supported;
// everything else is excluded from SupportedFeatures.
type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) List(request keeper.Request) ([]keeper.Response, int, error) {
	ctx, client, err := a.connect(request)
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
	for _, m := range memberList.Members {
		members = append(members, member{ID: m.ID, Name: m.Name, ClientURLs: m.ClientURLs, IsLearner: m.IsLearner})
		if len(m.ClientURLs) == 0 {
			continue
		}
		status, errStatus := client.Status(ctx, m.ClientURLs[0])
		if errStatus != nil {
			statuses[m.ID] = endpointStatus{Err: errStatus}
			continue
		}
		statuses[m.ID] = endpointStatus{Leader: status.Leader, RaftIndex: status.RaftIndex}
	}

	return mapMembers(members, statuses), http.StatusOK, nil
}

func (a *Adapter) Switchover(request keeper.Request) (*string, int, error) {
	body, errBody := parseSwitchoverBody(request.Body)
	if errBody != nil {
		return nil, http.StatusBadRequest, errBody
	}

	ctx, client, err := a.connect(request)
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

func (a *Adapter) Config(request keeper.Request) (any, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) ConfigUpdate(request keeper.Request) (any, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) DeleteSwitchover(request keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) Reinitialize(request keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) Restart(request keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) DeleteRestart(request keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) Reload(request keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) Failover(request keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) Activate(request keeper.Request) (*string, int, error) {
	return nil, http.StatusNotImplemented, keeper.ErrNotSupported
}

func (a *Adapter) Pause(request keeper.Request) (*string, int, error) {
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
	c.Client.Close()
}

// connect opens an etcd client for the request's endpoint and returns a
// context bound to the client's request timeout together with it.
func (a *Adapter) connect(request keeper.Request) (context.Context, *connection, error) {
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

func findLeader(statuses map[uint64]endpointStatus) (uint64, uint64) {
	var leaderID uint64
	for _, status := range statuses {
		if status.Err == nil && status.Leader != 0 {
			leaderID = status.Leader
			break
		}
	}
	if leader, ok := statuses[leaderID]; ok && leader.Err == nil {
		return leaderID, leader.RaftIndex
	}
	return leaderID, 0
}

func mapMember(m member, status endpointStatus, leaderID uint64, leaderRaftIndex uint64) keeper.Response {
	name := m.Name
	var keeperStatus keeper.Status = keeper.Active

	state := "running"
	var role keeper.Role = keeper.Unknown
	// NOTE: lag is the raft index difference to the leader, unlike patroni's lag value
	lag := int64(-1)
	switch {
	case len(m.ClientURLs) == 0:
		state = "starting"
	case status.Err != nil:
		state = "unreachable"
	case m.ID == leaderID:
		role = keeper.Leader
		lag = 0
	default:
		role = keeper.Replica
		if leaderRaftIndex >= status.RaftIndex {
			lag = int64(leaderRaftIndex - status.RaftIndex)
		}
	}

	var tags *map[string]any
	if m.IsLearner {
		tags = &map[string]any{"learner": true}
	}

	response := keeper.Response{
		Key:    &name,
		Status: &keeperStatus,
		State:  state,
		Role:   role,
		Lag:    lag,
		Tags:   tags,
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
