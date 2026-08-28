package node

import (
	"ivory/core/service/cert"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"

	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type KeeperConnection struct {
	Host string `json:"host" form:"host"`
	Port int    `json:"port" form:"port"`
}

// KeeperPlugin and KeeperStatus/KeeperRole/KeeperState are kept as aliases so
// request binding and keeperRegistry.Get() don't need explicit conversions
// at the call site.
type KeeperPlugin = keeper.Plugin

// PlatformPlugin selects which deployment target an operation runs against.
type PlatformPlugin = platform.Plugin

type KeeperStatus = keeper.Status
type KeeperRole = keeper.Role
type KeeperState = keeper.State

const (
	KeeperRoleLeader  KeeperRole = "leader"
	KeeperRoleUnknown KeeperRole = "unknown"
)

const KeeperStateUnknown = keeper.StateUnknown
const KeeperStateUnreachable = keeper.StateUnreachable

type KeeperScheduledSwitchover struct {
	At string `json:"at"`
	To string `json:"to"`
}

type KeeperScheduledRestart struct {
	At             string `json:"at"`
	PendingRestart bool   `json:"pendingRestart"`
}

type KeeperResponse struct {
	Key                  *string                    `json:"key"`
	Status               *KeeperStatus              `json:"status"`
	State                KeeperState                `json:"state"`
	Role                 KeeperRole                 `json:"role"`
	Sync                 bool                       `json:"sync"`
	Lag                  int64                      `json:"lag"`
	PendingRestart       bool                       `json:"pendingRestart"`
	ScheduledSwitchover  *KeeperScheduledSwitchover `json:"scheduledSwitchover"`
	ScheduledRestart     *KeeperScheduledRestart    `json:"scheduledRestart"`
	Tags                 *map[string]any            `json:"tags"`
	DiscoveredHost       *string                    `json:"discoveredHost"`
	DiscoveredName       *string                    `json:"discoveredName"`
	DiscoveredKeeperPort *int                       `json:"discoveredKeeperPort"`
	DiscoveredDbPort     *int                       `json:"discoveredDbPort"`
}

type KeeperOptions struct {
	Plugin  KeeperPlugin `json:"plugin" form:"plugin"`
	VaultId *uuid.UUID   `json:"vaultId" form:"vaultId"`
	Certs   *cert.Certs  `json:"certs" form:"certs"`
}

type KeeperOneRequest struct {
	KeeperConnection
	Body any `json:"body" form:"body"`
	KeeperOptions
}

type KeeperOneResponse = KeeperResponse

type KeeperMultiRequest struct {
	Connections []KeeperConnection `json:"connections" form:"connections"`
	Body        any                `json:"body" form:"body"`
	KeeperOptions
}

type KeeperMultiResponse struct {
	Connection KeeperConnection `json:"connection"`
	Response   []KeeperResponse `json:"response"`
	Error      string           `json:"error,omitempty"`
}

type CpuMetrics struct {
	TotalTicks uint64 `json:"totalTicks"`
	IdleTicks  uint64 `json:"idleTicks"`
}

type MemoryMetrics struct {
	TotalBytes     uint64 `json:"totalBytes"`
	AvailableBytes uint64 `json:"availableBytes"`
}

type NetworkMetrics struct {
	ReceivedBytes    uint64 `json:"receivedBytes"`
	TransmittedBytes uint64 `json:"transmittedBytes"`
}

type PlatformMetrics struct {
	Cpu     CpuMetrics     `json:"cpu"`
	Memory  MemoryMetrics  `json:"memory"`
	Network NetworkMetrics `json:"network"`
}

type PlatformVaultConnection struct {
	Host    string    `json:"host" form:"host" binding:"required"`
	Port    int       `json:"port" form:"port" binding:"required"`
	VaultId uuid.UUID `json:"vaultId" form:"vaultId" binding:"required"`
	// Platform selects the adapter; empty means Linux, so clusters stored
	// before platforms were selectable keep resolving.
	Platform PlatformPlugin `json:"platform" form:"platform"`
}

// PlatformOrDefault resolves the adapter key, defaulting to Linux.
func (c PlatformVaultConnection) PlatformOrDefault() PlatformPlugin {
	if c.Platform == "" {
		return platform.Linux
	}
	return c.Platform
}

type PlatformCredConnection struct {
	Host     string `json:"host" form:"host"`
	Port     int    `json:"port" form:"port"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type PlatformMetricsRequest = PlatformVaultConnection

type PlatformMetricsResponse = PlatformMetrics

type PlatformCopyIdRequest struct {
	PlatformCredConnection
	PublicKey string `json:"publicKey"`
}

// PlatformUpRequest is the low-level deployment primitive: Command is the
// user's own deployment command, still holding {{placeholder}} variables. The
// host comes from the connection and database credentials from the vault, so
// neither can be spoofed through Values.
type PlatformUpRequest struct {
	Connection PlatformVaultConnection `json:"connection" form:"connection" binding:"required"`
	Vaults     Vaults                  `json:"vaults" form:"vaults" binding:"required"`
	Command    string                  `json:"command" form:"command" binding:"required"`
	Values     keeper.Values           `json:"values" form:"values"`
}

type Vaults struct {
	// DatabaseId may be empty when the keeper plugin declares that it needs
	// no database credentials (Requirements.Credentials is false)
	DatabaseId uuid.UUID `json:"databaseId"`
	SshKeyId   uuid.UUID `json:"sshKeyId" binding:"required"`
}

// PlatformExecRequest runs a command inside the named deployment; Command is
// a template interpolated with Values plus the vault credentials, like
// PlatformUpRequest.
type PlatformExecRequest struct {
	Name       string                  `json:"name" binding:"required"`
	Connection PlatformVaultConnection `json:"connection" binding:"required"`
	Vaults     Vaults                  `json:"vaults" binding:"required"`
	Command    string                  `json:"command" binding:"required"`
	Values     keeper.Values           `json:"values"`
}

type KeeperDeploySpecRequest struct {
	Plugin KeeperPlugin `json:"plugin" form:"plugin" binding:"required"`
}

// KeeperDeploySpecResponse is plugin metadata the deploy forms need: the
// default endpoints and whether credentials are consumed. It says nothing
// about how to deploy - that is a command the user writes.
type KeeperDeploySpecResponse struct {
	DbPort      int    `json:"dbPort"`
	KeeperPort  *int   `json:"keeperPort"`
	Credentials bool   `json:"credentials"`
	DbUser      string `json:"dbUser"`
}

// KeeperDeployRequest deploys one node. It is deliberately flat: node owns no
// node struct of its own, and Host/ssh port come from Connection.
type KeeperDeployRequest struct {
	Plugin     KeeperPlugin            `json:"plugin" form:"plugin" binding:"required"`
	Cluster    string                  `json:"cluster" form:"cluster"`
	Name       string                  `json:"name" form:"name" binding:"required"`
	KeeperPort int                     `json:"keeperPort" form:"keeperPort"`
	DbPort     int                     `json:"dbPort" form:"dbPort"`
	Command    string                  `json:"command" form:"command" binding:"required"`
	PostScript string                  `json:"postScript" form:"postScript"`
	Connection PlatformVaultConnection `json:"connection" form:"connection" binding:"required"`
	Vaults     Vaults                  `json:"vaults" form:"vaults" binding:"required"`
}

type PlatformLogsRequest struct {
	Connection PlatformVaultConnection `json:"connection" form:"connection" binding:"required"`
	Path       string                  `json:"path" form:"path" binding:"required"`
	Tail       int                     `json:"tail" form:"tail"`
	Follow     bool                    `json:"follow" form:"follow"`
}

type PlatformActionRequest struct {
	Connection PlatformVaultConnection `json:"connection" form:"connection" binding:"required"`
	Name       string                  `json:"name" form:"name" binding:"required"`
}

type PlatformResponse = []string

type Process struct {
	Pid         int     `json:"pid"`
	Program     string  `json:"program"`
	Command     string  `json:"command"`
	Threads     int     `json:"threads"`
	User        string  `json:"user"`
	MemoryBytes uint64  `json:"memoryBytes"`
	MemPercent  float64 `json:"memPercent"`
	CpuPercent  float64 `json:"cpuPercent"`
}

type PlatformProcessesRequest = PlatformVaultConnection

type PlatformProcessesResponse = []Process

type InfoItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PlatformInfoRequest = PlatformVaultConnection

type PlatformInfoResponse = []InfoItem

// SPECIFIC (SERVER)

func mapKeeperResponse(r keeper.Response) KeeperResponse {
	var switchover *KeeperScheduledSwitchover
	if r.ScheduledSwitchover != nil {
		switchover = &KeeperScheduledSwitchover{At: r.ScheduledSwitchover.At, To: r.ScheduledSwitchover.To}
	}
	var restart *KeeperScheduledRestart
	if r.ScheduledRestart != nil {
		restart = &KeeperScheduledRestart{At: r.ScheduledRestart.At, PendingRestart: r.ScheduledRestart.PendingRestart}
	}
	return KeeperResponse{
		Key:                  r.Key,
		Status:               r.Status,
		State:                r.State,
		Role:                 r.Role,
		Sync:                 r.Sync,
		Lag:                  r.Lag,
		PendingRestart:       r.PendingRestart,
		ScheduledSwitchover:  switchover,
		ScheduledRestart:     restart,
		Tags:                 r.Tags,
		DiscoveredHost:       r.DiscoveredHost,
		DiscoveredName:       r.DiscoveredName,
		DiscoveredKeeperPort: r.DiscoveredKeeperPort,
		DiscoveredDbPort:     r.DiscoveredDbPort,
	}
}

// mapKeeperDeploymentToPlatformSpec converts keeper deployment requirements
// into a platform deploy spec. Name, hostname and restart policy are Ivory
// deployment policy rather than keeper knowledge. Single-host mode uses host
// networking and drops port mappings, volumes and the restart policy because
// it targets development and testing setups.

func mapPlatformMetrics(m *platform.Metrics) *PlatformMetrics {
	if m == nil {
		return nil
	}
	return &PlatformMetrics{
		Cpu:     CpuMetrics{TotalTicks: m.Cpu.TotalTicks, IdleTicks: m.Cpu.IdleTicks},
		Memory:  MemoryMetrics{TotalBytes: m.Memory.TotalBytes, AvailableBytes: m.Memory.AvailableBytes},
		Network: NetworkMetrics{ReceivedBytes: m.Network.ReceivedBytes, TransmittedBytes: m.Network.TransmittedBytes},
	}
}

func mapPlatformProcesses(processes []platform.Process) PlatformProcessesResponse {
	response := make(PlatformProcessesResponse, 0, len(processes))
	for _, p := range processes {
		response = append(response, Process{
			Pid:         p.Pid,
			Program:     p.Program,
			Command:     p.Command,
			Threads:     p.Threads,
			User:        p.User,
			MemoryBytes: p.MemoryBytes,
			MemPercent:  p.MemPercent,
			CpuPercent:  p.CpuPercent,
		})
	}
	return response
}

func mapPlatformInfo(items []platform.InfoItem) PlatformInfoResponse {
	response := make(PlatformInfoResponse, 0, len(items))
	for _, i := range items {
		response = append(response, InfoItem{Key: i.Key, Value: i.Value})
	}
	return response
}
