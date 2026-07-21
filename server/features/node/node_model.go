package node

import (
	"ivory/core/service/cert"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
	"maps"

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

// Built-in interpolation variable names, re-exported so other features never
// import the plugin package directly.
const (
	VarCluster    = keeper.VarCluster
	VarHost       = keeper.VarHost
	VarKeeperPort = keeper.VarKeeperPort
	VarDbPort     = keeper.VarDbPort
	VarDbUser     = keeper.VarDbUser
	VarDbPass     = keeper.VarDbPass
)

type KeeperStatus = keeper.Status
type KeeperRole = keeper.Role
type KeeperState = keeper.State

const (
	KeeperRoleLeader  KeeperRole = "leader"
	KeeperRoleUnknown KeeperRole = "unknown"
)

const KeeperStateUnknown = keeper.StateUnknown

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

// PlatformUpRequest is the low-level deployment primitive: Options is the
// user-editable options template and Values holds the {{placeholder}}
// interpolation values (cluster, dcs, ports, aux ports, ...). The host is
// taken from the connection and database credentials from the vault, so they
// cannot be spoofed through Values.
type PlatformUpRequest struct {
	Name       string                  `json:"name" form:"name" binding:"required"`
	Image      string                  `json:"image" form:"image" binding:"required"`
	Connection PlatformVaultConnection `json:"connection" form:"connection" binding:"required"`
	Vaults     Vaults                  `json:"vaults" form:"vaults" binding:"required"`
	Options    string                  `json:"options" form:"options" binding:"required"`
	Values     map[string]string       `json:"values" form:"values"`
}

type Vaults struct {
	// DatabaseId may be empty when the keeper plugin declares that it needs
	// no database credentials (DeployFieldsResponse.Defaults has no {{dbUser}})
	DatabaseId uuid.UUID `json:"databaseId"`
	SshKeyId   uuid.UUID `json:"sshKeyId" binding:"required"`
}

// PlatformExecRequest runs a command inside the named deployment; Command is
// a template interpolated with Values plus the vault credentials, like
// PlatformUpRequest options.
type PlatformExecRequest struct {
	Name       string                  `json:"name" binding:"required"`
	Connection PlatformVaultConnection `json:"connection" binding:"required"`
	Vaults     Vaults                  `json:"vaults" binding:"required"`
	Command    string                  `json:"command" binding:"required"`
	Values     map[string]string       `json:"values"`
}

type KeeperDeploySpecRequest struct {
	Plugin KeeperPlugin `json:"plugin" form:"plugin" binding:"required"`
}

type KeeperDeploySpecResponse struct {
	Uri    string               `json:"uri"`
	Fields DeployFieldsResponse `json:"fields"`
}

// KeeperDeployPlanRequest describes a deployment intent: everything except
// the node hosts is optional and falls back to the keeper plugin's
// DeploymentSpec. Values carries plugin-required inputs (e.g. "dcs" for
// plugins with a manual DCS) and extra interpolation values.
type KeeperDeployPlanRequest struct {
	Plugin     KeeperPlugin                  `json:"plugin" form:"plugin" binding:"required"`
	Cluster    string                        `json:"cluster" form:"cluster"`
	SingleHost bool                          `json:"singleHost" form:"singleHost"`
	Image      string                        `json:"image" form:"image"`
	Values     map[string]string             `json:"values" form:"values"`
	Nodes      []KeeperDeployPlanNodeRequest `json:"nodes" form:"nodes"`
}

type KeeperDeployPlanNodeRequest struct {
	Host       string `json:"host" binding:"required"`
	SshPort    *int   `json:"sshPort"`
	KeeperPort *int   `json:"keeperPort"`
	DbPort     *int   `json:"dbPort"`
	// Options overrides the rendered options template for this node.
	Options string `json:"options"`
}

// KeeperDeployPlanResponse is the resolved deployment: concrete ports and
// options per node, the effective field values (user-provided or computed),
// the post-deploy command templates, and advisory warnings (missing
// placeholder values, ignored ports). Previews mask credentials.
type KeeperDeployPlanResponse struct {
	Image      string                 `json:"image"`
	Values     map[string]string      `json:"values"`
	PostDeploy []string               `json:"postDeploy"`
	Fields     DeployFieldsResponse   `json:"fields"`
	Nodes      []KeeperDeployPlanNode `json:"nodes"`
	Warnings   []string               `json:"warnings"`
}

type KeeperDeployPlanNode struct {
	Host       string         `json:"host"`
	SshPort    int            `json:"sshPort"`
	KeeperPort int            `json:"keeperPort"`
	DbPort     int            `json:"dbPort"`
	Ports      map[string]int `json:"ports"`
	Options    string         `json:"options"`
	Preview    string         `json:"preview"`
}

// DeployFieldResponse describes one editable image-level field: its value
// interpolates as {{name}}, the plan prefills it (user edit wins), Derived
// marks values computed from the node list.
type DeployFieldResponse struct {
	Name    string `json:"name"`
	Label   string `json:"label"`
	Example string `json:"example,omitempty"`
	Type    string `json:"type"`
	Default string `json:"default,omitempty"`
	Derived bool   `json:"derived"`
}

// DeployFieldsResponse tells the frontend which deploy form fields the keeper
// plugin needs. Defaults mirrors the spec's built-in variable defaults keyed
// by the variable's interpolated form: an absent {{keeperPort}} hides the
// keeper port inputs (the keeper endpoint is the database itself), an absent
// {{dbUser}} hides the credential inputs and drops the vault requirement (a
// non-empty value is the engine-required username, prefilled and locked).
// Fields render as editable prefilled inputs.
type DeployFieldsResponse struct {
	Defaults map[string]string     `json:"defaults"`
	Fields   []DeployFieldResponse `json:"fields"`
}

// KeeperDeployRequest deploys a single keeper node end-to-end: Plugin and
// Values resolve the deployment plan for this one node (ports, options,
// interpolation), Connection and Vaults are resolved by the caller (a stored
// cluster's vaults or a freshly entered SSH/database credential) since node
// has no access to the cluster feature's storage.
type KeeperDeployRequest struct {
	Plugin     KeeperPlugin                `json:"plugin" form:"plugin" binding:"required"`
	Cluster    string                      `json:"cluster" form:"cluster"`
	Image      string                      `json:"image" form:"image"`
	Values     map[string]string           `json:"values" form:"values"`
	Node       KeeperDeployPlanNodeRequest `json:"node" form:"node"`
	Connection PlatformVaultConnection     `json:"connection" form:"connection" binding:"required"`
	Vaults     Vaults                      `json:"vaults" form:"vaults" binding:"required"`
}

// KeeperDeployUpRequest deploys one node that a KeeperDeployPlan already
// resolved: PlanValues are the plan's effective field values (dcs, derived
// member lists, ...), RequestValues the raw request-supplied interpolation
// extras, and Node the node's own resolved ports/options.
type KeeperDeployUpRequest struct {
	Cluster       string
	Image         string
	PlanValues    map[string]string
	RequestValues map[string]string
	Node          KeeperDeployPlanNode
	Connection    PlatformVaultConnection
	Vaults        Vaults
}

// KeeperPostDeployRequest runs a deployment plan's post-deploy commands (e.g.
// enabling authentication) inside one already-deployed node.
type KeeperPostDeployRequest struct {
	Cluster       string
	RequestValues map[string]string
	PlanValues    map[string]string
	PostDeploy    []string
	Node          KeeperDeployPlanNode
	Connection    PlatformVaultConnection
	Vaults        Vaults
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
		DiscoveredKeeperPort: r.DiscoveredKeeperPort,
		DiscoveredDbPort:     r.DiscoveredDbPort,
	}
}

// mapKeeperDeploymentToPlatformSpec converts keeper deployment requirements
// into a platform deploy spec. Name, hostname and restart policy are Ivory
// deployment policy rather than keeper knowledge. Single-host mode uses host
// networking and drops port mappings, volumes and the restart policy because
// it targets development and testing setups.
func mapKeeperDeploymentToPlatformSpec(s keeper.DeploymentSpec, singleHost bool) platform.DeploySpec {
	spec := platform.DeploySpec{
		Name:        keeper.VarHost,
		Hostname:    keeper.VarHost,
		HostNetwork: singleHost,
	}
	if !singleHost {
		spec.RestartPolicy = "unless-stopped"
		spec.Ports = s.Ports
		for _, v := range s.Volumes {
			spec.Volumes = append(spec.Volumes, platform.VolumeMount{HostPath: v.HostPath, ContainerPath: v.ContainerPath})
		}
	}
	for _, e := range s.Env {
		spec.Env = append(spec.Env, platform.EnvVar{Name: e.Name, Value: e.Value})
	}
	return spec
}

func mapKeeperDeploymentFields(s keeper.DeploymentSpec) DeployFieldsResponse {
	defaults := make(map[string]string, len(s.Defaults))
	maps.Copy(defaults, s.Defaults)
	fields := make([]DeployFieldResponse, 0, len(s.Fields))
	for _, f := range s.Fields {
		fields = append(fields, DeployFieldResponse{
			Name:    f.Name,
			Label:   f.Label,
			Example: f.Example,
			Type:    string(f.Type),
			Default: f.Default,
			Derived: f.Template != "",
		})
	}
	return DeployFieldsResponse{
		Defaults: defaults,
		Fields:   fields,
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
