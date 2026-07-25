import {Certs} from "../../cert/api/CertType"

// COMMON (WEB AND SERVER)

export enum KeeperPlugin {
    PATRONI_POSTGRES = "patroni_postgres",
    NATIVE_POSTGRES = "native_postgres",
    NATIVE_ETCD = "native_etcd",
}

export interface KeeperConnection {
    host: string,
    port: number,
}

export enum KeeperStatus {
    Active = "ACTIVE",
    Paused = "PAUSED",
}

export type Role = "leader"|"replica"|"unknown";

// KeeperState is Ivory's normalized keeper state: every adapter maps its own
// plugin-specific and version-specific vocabulary (e.g. Patroni renamed
// state names across releases) onto this fixed set on the backend.
export type KeeperState = "running"|"starting"|"restarting"|"stopping"|"stopped"|"failed"|"unreachable"|"unknown";

export interface Keeper {
    host: string,
    port: number,
    name?: string,
    status?: KeeperStatus,
}

export interface KeeperResponse {
    key?: string,
    status?: KeeperStatus,
    state: KeeperState,
    role: Role,
    sync: boolean,
    lag: number,
    pendingRestart: boolean,
    scheduledSwitchover?: ScheduledSwitchover,
    scheduledRestart?: ScheduledRestart,
    tags?: {[key: string]: any},
    discoveredHost?: string,
    discoveredKeeperPort?: number,
    discoveredDbPort?: number,
}

export interface ScheduledSwitchover {
    at: string,
    to: string,
}

export interface ScheduledRestart {
    at: string,
    pendingRestart: boolean,
}

export interface KeeperOneRequest extends KeeperConnection {
    plugin: KeeperPlugin,
    vaultId?: string,
    certs?: Certs,
    body?: any,
}

export type KeeperOneResponse = KeeperResponse

export interface KeeperMultiRequest {
    connections: KeeperConnection[],
    body?: any,
}

export interface KeeperMultiResponse {
    connection: KeeperConnection,
    response: KeeperResponse[],
    error?: string,
}

export interface PlatformVaultConnection {
    host: string,
    port: number,
    vaultId: string,
}

export interface PlatformCredConnection {
    host: string,
    port: number,
    username: string,
    password?: string,
}

export interface CpuMetrics {
    totalTicks: number,
    idleTicks: number,
}

export interface MemoryMetrics {
    totalBytes: number,
    availableBytes: number,
}

export interface NetworkMetrics {
    receivedBytes: number,
    transmittedBytes: number,
}

export interface Metrics {
    cpu: CpuMetrics,
    memory: MemoryMetrics,
    network: NetworkMetrics,
}

export type PlatformMetricsResponse = Metrics

export type PlatformMetricsRequest = PlatformVaultConnection

export interface PlatformCopyIdRequest {
    host: string,
    port: number,
    username: string,
    password: string,
    publicKey: string,
}

// PlatformUpRequest is the low-level deployment primitive: options is the
// user-editable options template and values holds the {{placeholder}}
// interpolation values (cluster, dcs, ports, aux ports, ...).
export interface PlatformUpRequest {
    name: string,
    image: string,
    connection: PlatformVaultConnection,
    vaults: {
        // NOTE: optional for keeper plugins that consume no database
        // credentials (no {{dbUser}} in the keeper plugin's spec defaults)
        databaseId?: string,
        sshKeyId: string,
    },
    options: string,
    values: {[key: string]: string},
}

// KeeperDeployNode is the single-node shape for KeeperDeployRequest: it
// mirrors the server's KeeperDeployPlanNodeRequest rather than reusing the
// cluster feature's DeployNode, since node must not import from cluster.
export interface KeeperDeployNode {
    host: string,
    sshPort?: number,
    keeperPort?: number,
    dbPort?: number,
    // options overrides the rendered options template for this node
    options?: string,
}

// KeeperDeployRequest deploys a single keeper node end-to-end: plugin and
// values resolve the deployment plan for this one node (ports, options,
// interpolation), connection and vaults are resolved by the caller (e.g. a
// stored cluster's vaults) since node has no access to cluster storage.
export interface KeeperDeployRequest {
    plugin: KeeperPlugin,
    cluster: string,
    image?: string,
    values: {[key: string]: string},
    node: KeeperDeployNode,
    connection: PlatformVaultConnection,
    vaults: {
        databaseId?: string,
        sshKeyId: string,
    },
}

export interface KeeperDeploySpecRequest {
    plugin: KeeperPlugin,
}

// DeployFieldResponse describes one editable image-level field: its value
// interpolates as {{name}}, the plan prefills it (a user edit wins), derived
// marks values computed from the node list.
export interface DeployFieldResponse {
    name: string,
    label: string,
    example?: string,
    type: "text" | "port",
    default?: string,
    derived: boolean,
}

// InterpolationVar is a built-in {{placeholder}} variable usable in the
// user-editable Options text, in its interpolated form; plugin-declared
// field names extend this set. Mirrors the subset of the server's
// keeper.Var constants that can appear in Options - keeper.VarPrimaryHost
// is deliberately not included here, since it only ever appears in
// EntryScript, which is never shown or editable in the UI.
export enum InterpolationVar {
    Cluster = "{{cluster}}",
    Host = "{{host}}",
    KeeperPort = "{{keeperPort}}",
    DbPort = "{{dbPort}}",
    DbUser = "{{dbUser}}",
    DbPass = "{{dbPass}}",
}

// DeployFieldsResponse tells the form which fields the keeper plugin needs.
// defaults mirrors the spec's built-in variable defaults keyed by the
// variable's interpolated form: an absent {{keeperPort}} hides the keeper
// port inputs (the keeper endpoint is the database itself), an absent
// {{dbUser}} hides the credential inputs (a non-empty value is the
// engine-required username, prefilled and locked).
export interface DeployFieldsResponse {
    defaults: {[name: string]: string},
    fields: DeployFieldResponse[],
}

export interface KeeperDeploySpecResponse {
    uri: string,
    fields: DeployFieldsResponse,
}

// KeeperDeployPlanRequest describes a deployment intent: everything except the
// node hosts is optional and falls back to the keeper plugin's spec.
export interface KeeperDeployPlanRequest {
    plugin: KeeperPlugin,
    cluster: string,
    singleHost: boolean,
    image?: string,
    values: {[key: string]: string},
    nodes: KeeperDeployNode[],
}

// KeeperDeployPlanResponse is the resolved deployment: concrete ports and options
// per node, the effective field values (user-provided or computed), the
// post-deploy command templates, and advisory warnings. Previews mask
// credentials.
export interface KeeperDeployPlanResponse {
    image: string,
    values: {[name: string]: string},
    postDeploy: string[],
    fields: DeployFieldsResponse,
    nodes: KeeperDeployPlanNode[],
    warnings: string[],
}

export interface KeeperDeployPlanNode {
    host: string,
    sshPort: number,
    keeperPort: number,
    dbPort: number,
    ports: {[name: string]: number},
    options: string,
    optionsPreview: string,
    entryScript: string,
    entryScriptPreview: string,
}

export interface PlatformLogsRequest {
    connection: PlatformVaultConnection,
    path: string,
    tail?: number,
    follow?: boolean,
}

export interface PlatformActionRequest {
    connection: PlatformVaultConnection,
    name: string,
}

export interface Process {
    pid: number,
    program: string,
    command: string,
    threads: number,
    user: string,
    memoryBytes: number,
    memPercent: number,
    cpuPercent: number,
}

export type PlatformProcessesResponse = Process[]

export interface InfoItem {
    key: string,
    value: string,
}

export type PlatformInfoResponse = InfoItem[]

export type PlatformInfoRequest = PlatformVaultConnection

// SPECIFIC (WEB)

export enum NodeTabType {DATABASE, CONTAINER, KEEPER, TOOLS, PLATFORM}
