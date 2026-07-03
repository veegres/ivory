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

// KeeperPlugin and KeeperStatus/KeeperRole are kept as aliases so request binding and
// keeperRegistry.Get() don't need explicit conversions at the call site.
type KeeperPlugin = keeper.Plugin
type KeeperStatus = keeper.Status
type KeeperRole = keeper.Role

const (
	KeeperRoleLeader  KeeperRole = "leader"
	KeeperRoleUnknown KeeperRole = "unknown"
)

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
	State                string                     `json:"state"`
	Role                 KeeperRole                 `json:"role"`
	Lag                  int64                      `json:"lag"`
	PendingRestart       bool                       `json:"pendingRestart"`
	ScheduledSwitchover  *KeeperScheduledSwitchover `json:"scheduledSwitchover"`
	ScheduledRestart     *KeeperScheduledRestart    `json:"scheduledRestart"`
	Tags                 *map[string]any            `json:"tags"`
	DiscoveredHost       *string                    `json:"discoveredHost"`
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

type PlatformUpRequest struct {
	Name            string                  `json:"name" form:"name" binding:"required"`
	Image           string                  `json:"image" form:"image" binding:"required"`
	Connection      PlatformVaultConnection `json:"connection" form:"connection" binding:"required"`
	Vaults          Vaults                  `json:"vaults" form:"vaults" binding:"required"`
	ImageOptions    ImageOptionsRequest     `json:"imageOptions" form:"imageOptions" binding:"required"`
	RawImageOptions string                  `json:"rawImageOptions" form:"rawImageOptions" binding:"required"`
}

type Vaults struct {
	DatabaseId uuid.UUID `json:"databaseId" binding:"required"`
	SshKeyId   uuid.UUID `json:"sshKeyId" binding:"required"`
}

type ImageOptionsRequest struct {
	Cluster    string `json:"cluster" binding:"required"`
	Dcs        string `json:"dcs" binding:"required"`
	KeeperPort int    `json:"keeperPort" binding:"required"`
	DbPort     int    `json:"dbPort" binding:"required"`
}

type ImageOptions struct {
	Host       string `json:"host"`
	Cluster    string `json:"cluster"`
	Dcs        string `json:"dcs"`
	DbPass     string `json:"dbPass"`
	DbUser     string `json:"dbUser"`
	KeeperPort string `json:"keeperPort"`
	DbPort     string `json:"dbPort"`
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
		Lag:                  r.Lag,
		PendingRestart:       r.PendingRestart,
		ScheduledSwitchover:  switchover,
		ScheduledRestart:     restart,
		Tags:                 r.Tags,
		DiscoveredHost:       r.DiscoveredHost,
		DiscoveredKeeperPort: r.DiscoveredKeeperPort,
		DiscoveredDbPort:     r.DiscoveredDbPort,
	}
}

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
