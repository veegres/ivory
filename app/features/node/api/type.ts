import {ReactNode} from "react"

import {Certs} from "../../cert/api/type"
import {Cluster, Node} from "../../cluster/api/type"
import {Feature} from "../../feature"

// COMMON (WEB AND SERVER)

export interface KeeperConnection {
    host: string,
    port: number,
}

export enum KeeperPlugin {
    PATRONI = "patroni",
    POSTGRES = "postgres",
}

export enum KeeperStatus {
    Active = "ACTIVE",
    Paused = "PAUSED",
}

export type Role = "leader"|"replica"|"unknown";

export interface Keeper {
    host: string,
    port: number,
    name?: string,
    status?: KeeperStatus,
}

export interface KeeperResponse {
    key?: string,
    status?: KeeperStatus,
    state: string,
    role: Role,
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

export interface PlatformConnection {
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

export interface PlatformMetricsRequest {
    host: string,
    port: number,
    vaultId: string,
}

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
    connection: PlatformConnection,
    vaults: {
        databaseId: string,
        sshKeyId: string,
    },
    imageOptions: ImageOptions,
    rawImageOptions: string,
}

export interface ImageOptions {
    cluster: string,
    dcs: string,
    host: string,
    keeperPort: number,
    dbPort: number,
}

export interface PlatformLogsRequest {
    connection: PlatformConnection,
    path: string,
    tail?: number,
    follow?: boolean,
}

export interface PlatformActionRequest {
    connection: PlatformConnection,
    name: string,
}

// SPECIFIC (WEB)

export enum NodeTabType {DATABASE, CONTAINER, KEEPER, TOOLS, PLATFORM}
export interface NodeTab {
    label: string,
    feature?: Feature,
    body: (cluster: Cluster, node: Node) => ReactNode,
    info?: ReactNode,
    actions?: ReactNode,
}
