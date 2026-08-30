package keeper

import (
	"crypto/tls"
)

// COMMON (WEB AND SERVER)

// PluginType identifies which keeper (HA management tool) manages a node.
type PluginType string

const (
	PATRONI_POSTGRES  PluginType = "patroni_postgres"
	NATIVE_POSTGRES   PluginType = "native_postgres"
	NATIVE_ETCD       PluginType = "native_etcd"
	NATIVE_REDIS      PluginType = "native_redis"
	NATIVE_CLICKHOUSE PluginType = "native_clickhouse"
	NATIVE_ZOOKEEPER  PluginType = "native_zookeeper"
	NATIVE_MONGO      PluginType = "native_mongo"
)

func (p PluginType) String() string {
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
	Replica Role = "replica"
	Unknown Role = "unknown"
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
	// Key is the keeper's own identifier for this member (e.g. a Patroni
	// member name, or host:port for native postgres, which has no separate
	// member name). It is opaque outside the keeper - callers use it to
	// refer back to this exact member (e.g. as the current leader's key
	// when requesting a switchover), not to parse or derive anything from it.
	Key *string `json:"key"`
	// Status is this member's keeper-wide Active/Paused flag (see the
	// Status type doc). It reflects the whole keeper's failover management,
	// not anything specific to this one member.
	Status *Status `json:"status"`
	// State is this member's current State (see the State type doc).
	State State `json:"state"`
	// Role is this member's current Role (see the Role type doc).
	Role Role `json:"role"`
	// Sync reports whether a Replica currently belongs to the keeper's
	// synchronous replication set, as opposed to a plain asynchronous
	// replica. Always false for Leader/Unknown and for keepers (e.g. etcd)
	// that have no synchronous-replica concept.
	Sync bool `json:"sync"`
	// Lag is this member's replication lag. It is only meaningful for a
	// Replica - a Leader always reports 0. The unit and exact measurement
	// differ per adapter (e.g. Patroni's own /cluster lag value vs. native
	// postgres' pg_wal_lsn_diff of received/replayed WAL), so Lag should
	// only be compared within one keeper plugin, never across plugins.
	Lag int64 `json:"lag"`
	// PendingRestart reports whether the keeper has flagged this member as
	// needing a restart (typically because a pending config change requires
	// one to take effect).
	PendingRestart bool `json:"pendingRestart"`
	// ScheduledSwitchover is set only on the member a pending switchover
	// would move the leader away from (see the ScheduledSwitchover doc);
	// nil for every other member.
	ScheduledSwitchover *ScheduledSwitchover `json:"scheduledSwitchover"`
	// ScheduledRestart is this member's own pending restart schedule, if
	// any (see the ScheduledRestart doc).
	ScheduledRestart *ScheduledRestart `json:"scheduledRestart"`
	// Tags is an opaque passthrough of whatever keeper-specific tags this
	// member reports (e.g. Patroni member tags), surfaced to the UI as-is
	// without Ivory interpreting their meaning.
	Tags *map[string]any `json:"tags"`

	// Discovered Topology (Crucial for Auto-Creation)
	//
	// These three fields describe a node's real connection details as
	// reported by the KEEPER ITSELF - e.g. Patroni's /cluster gives every
	// member's host/port parsed from its own api_url; native postgres's
	// primary can read a connected standby's address from
	// pg_stat_replication. They are ground truth independent of whatever
	// Ivory already has configured, which is exactly what makes them useful
	// for two things: building a cluster from scratch by asking a live
	// keeper (cluster.Detect), and detecting drift between what's
	// configured and what the keeper actually reports
	// (cluster.addOverviewWarnings' "keeper response and cluster
	// configuration mismatch" check).
	//
	// Because of that, an adapter must never populate these by copying back
	// Ivory's own already-configured NodeConfig - that isn't discovery,
	// it's echoing configuration, and it would silently defeat the drift
	// check above. If an adapter cannot determine a value from the keeper's
	// own response, leave it nil rather than guessing or borrowing it from
	// elsewhere.
	DiscoveredHost *string `json:"discoveredHost"`
	// DiscoveredName is the member's own name as the keeper itself calls it
	// (e.g. a Patroni member name, an etcd member name). It is the name the
	// node was deployed under, which is why it becomes NodeConfig.Name on
	// auto-creation instead of the host. An adapter whose engine has no
	// separate member name - one that identifies members by host:port
	// (mongo, redis, clickhouse, zookeeper, native postgres) - leaves it nil
	// rather than reporting an endpoint as a name.
	DiscoveredName *string `json:"discoveredName"`
	// DiscoveredKeeperPort is required whenever DiscoveredHost is set: a
	// host known without a port can't be matched back to a configured node
	// (see cluster.addKeeperResponsesToMap), so an adapter that discovers a
	// peer's host must also resolve its port from data it itself possesses
	// (e.g. native postgres's Adapter uses its own connection port, since
	// every node in a native postgres cluster shares one port by
	// convention - see Adapter's doc comment) rather than leaving it unset.
	DiscoveredKeeperPort *int `json:"discoveredKeeperPort"`
	// DiscoveredDbPort is the node's real database port, following the same
	// rules as DiscoveredKeeperPort above (nil if the adapter can't
	// determine it from the keeper's own response).
	DiscoveredDbPort *int `json:"discoveredDbPort"`
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
