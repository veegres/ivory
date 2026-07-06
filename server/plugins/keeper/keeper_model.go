package keeper

import (
	"crypto/tls"
)

// COMMON (WEB AND SERVER)

// Plugin identifies which keeper (HA management tool) manages a node.
type Plugin string

const (
	PATRONI_POSTGRES Plugin = "patroni_postgres"
	NATIVE_POSTGRES  Plugin = "native_postgres"
	NATIVE_ETCD      Plugin = "native_etcd"
)

func (p Plugin) String() string {
	return string(p)
}

// Status reports whether a keeper is actively managing failover/switchover
// or has been Paused (see Adapter.Pause/Activate).
type Status string

const (
	Paused Status = "PAUSED"
	Active        = "ACTIVE"
)

// Role is a member's position within the cluster as seen by its keeper.
type Role string

const (
	Leader  Role = "leader"
	Replica      = "replica"
	Unknown      = "unknown"
)

// State is Ivory's normalized view of what the keeper's postgres process is
// doing right now. Every adapter maps its own plugin-specific and
// version-specific vocabulary onto this fixed set, so callers (the cluster
// overview, the frontend) only ever need to understand these values.
type State string

const (
	StateRunning     State = "running"
	StateStarting    State = "starting"
	StateRestarting  State = "restarting"
	StateStopping    State = "stopping"
	StateStopped     State = "stopped"
	StateFailed      State = "failed"
	StateUnreachable State = "unreachable"
	StateUnknown     State = "unknown"
)

// Keeper is the persisted, user-facing keeper connection configuration for a
// node.
type Keeper struct {
	Host   string  `json:"host"`
	Port   int     `json:"port"`
	Name   *string `json:"name"`
	Status *Status `json:"status"`
}

// Response is a single cluster member as reported by Adapter.List, combining
// the keeper's own state with database/keeper endpoints it discovered so
// Ivory can auto-create the corresponding node.
type Response struct {
	Key                 *string              `json:"key"`
	Status              *Status              `json:"status"`
	State               State                `json:"state"`
	Role                Role                 `json:"role"`
	Lag                 int64                `json:"lag"`
	PendingRestart      bool                 `json:"pendingRestart"`
	ScheduledSwitchover *ScheduledSwitchover `json:"scheduledSwitchover"`
	ScheduledRestart    *ScheduledRestart    `json:"scheduledRestart"`
	Tags                *map[string]any      `json:"tags"`

	// Discovered Topology (Crucial for Auto-Creation)
	DiscoveredHost       *string `json:"discoveredHost"`
	DiscoveredKeeperPort *int    `json:"discoveredKeeperPort"`
	DiscoveredDbPort     *int    `json:"discoveredDbPort"`
}

// ScheduledSwitchover describes a pending, not-yet-performed switchover.
type ScheduledSwitchover struct {
	At string `json:"at"`
	To string `json:"to"`
}

// ScheduledRestart describes a pending, not-yet-performed restart.
type ScheduledRestart struct {
	At             string `json:"at"`
	PendingRestart bool   `json:"pendingRestart"`
}

// SPECIFIC (SERVER)

// Request is the input to every Adapter method: the keeper's endpoint,
// optional credentials/TLS to reach it, and an operation-specific body.
type Request struct {
	Host        string       `json:"host" form:"host"`
	Port        int          `json:"port" form:"port"`
	Credentials *Credentials `json:"credentials" form:"credentials"`
	TlsConfig   *tls.Config  `json:"tlsConfig" form:"tlsConfig"`
	Body        any          `json:"body" form:"body"`
}

// Credentials authenticates a Request against its keeper.
type Credentials struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
