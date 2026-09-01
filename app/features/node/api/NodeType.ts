import {Certs} from "../../cert/api/CertType"

// COMMON (WEB AND SERVER)

export enum KeeperPlugin {
    PATRONI_POSTGRES = "patroni_postgres",
    NATIVE_POSTGRES = "native_postgres",
    NATIVE_ETCD = "native_etcd",
    NATIVE_REDIS = "native_redis",
    NATIVE_CLICKHOUSE = "native_clickhouse",
    NATIVE_ZOOKEEPER = "native_zookeeper",
    NATIVE_MONGO = "native_mongo",
}

// PlatformPlugin selects the deployment target a template and a deployment
// run against. Only docker exists today; k8s is the reason it is a choice.
export enum PlatformPlugin {
    DOCKER = "docker",
}

// DeployVar is the closed set of {{variables}} a deployment command may use:
// the node's own identity and endpoints, plus the credentials resolved from the
// vault and the coordination store the cluster runs against. Anything else an
// engine needs is written literally into the command. A placeholder outside
// this set is a validation error, never a new variable.
export enum DeployVar {
    Cluster = "{{cluster}}",
    Dcs = "{{dcs}}",
    Name = "{{name}}",
    Host = "{{host}}",
    SshPort = "{{sshPort}}",
    KeeperPort = "{{keeperPort}}",
    DbPort = "{{dbPort}}",
    KeeperUser = "{{keeperUser}}",
    KeeperPass = "{{keeperPass}}",
    DbUser = "{{dbUser}}",
    DbPass = "{{dbPass}}",
}

// DeployVarScope is how many values a variable has in one deployment: a Node
// variable is answered once per node, a Cluster variable once for the whole
// deployment and reaches every node's command with the same value. It is what
// the template editor groups the variable list by, and it says where the deploy
// screen asks for each one - the node card, or the sections above it.
export enum DeployVarScope {
    Command = "node",
    Template = "cluster",
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
    discoveredName?: string,
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
    // platform selects the deployment target; omitted means docker, so clusters
    // stored before platforms were selectable keep resolving
    platform?: PlatformPlugin,
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
    platform: PlatformPlugin,
    publicKey: string,
}

export interface PlatformUpRequest {
    connection: PlatformVaultConnection,
    vaults: DeployVaults,
    command: string,
}

export interface DeployVaults {
    keeperId?: string,
    databaseId?: string,
    sshKeyId: string,
}

export interface KeeperDeployRequest {
    plugin: KeeperPlugin,
    cluster: string,
    dcs?: string,
    name: string,
    keeperPort: number,
    dbPort: number,
    command: string,
    postScripts?: string[],
    connection: PlatformVaultConnection,
    vaults: DeployVaults,
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

export enum NodeTabType {DATABASE, CONTAINER, KEEPER, TOOLS, SYSTEM}

export enum ReleaseStage {
    ALPHA = "alpha",
    BETA = "beta",
    STABLE = "stable",
}
