import {Certs} from "../../cert/api/CertType"

// COMMON (WEB AND SERVER)

export enum KeeperPlugin {
    PATRONI_POSTGRES = "patroni_postgres",
    NATIVE_POSTGRES = "native_postgres",
    NATIVE_ETCD = "native_etcd",
    NATIVE_REDIS = "native_redis",
    NATIVE_CLICKHOUSE = "native_clickhouse",
    NATIVE_ZOOKEEPER = "native_zookeeper",
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

export interface KeeperDeployNode {
    host: string,
    sshPort?: number,
    keeperPort?: number,
    dbPort?: number,
    // options overrides the rendered options template for this node
    options?: string,
}

export interface KeeperDeployRequest {
    plugin: KeeperPlugin,
    cluster: string,
    singleHost: boolean,
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

export interface DeployFieldResponse {
    name: string,
    label: string,
    example?: string,
    type: "text" | "port",
    default?: string,
    derived: boolean,
}

export enum InterpolationVar {
    Cluster = "{{cluster}}",
    Host = "{{host}}",
    KeeperPort = "{{keeperPort}}",
    DbPort = "{{dbPort}}",
    DbUser = "{{dbUser}}",
    DbPass = "{{dbPass}}",
}

export interface DeployFieldsResponse {
    defaults: {[name: string]: string},
    fields: DeployFieldResponse[],
}

export interface KeeperDeploySpecResponse {
    uri: string,
    fields: DeployFieldsResponse,
}

export interface KeeperDeployPlanRequest {
    plugin: KeeperPlugin,
    cluster: string,
    singleHost: boolean,
    image?: string,
    values: {[key: string]: string},
    nodes: KeeperDeployNode[],
}

export interface KeeperDeployPlanResponse {
    image: string,
    values: {[name: string]: string},
    postScript: string,
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

export enum ReleaseStage {
    ALPHA = "alpha",
    BETA = "beta",
    STABLE = "stable",
}
