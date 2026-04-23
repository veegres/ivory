
import {ReactNode} from "react"

import {Certs} from "../cert/type"
import {Connection as QueryConnection} from "../query/type"

// COMMON (WEB AND SERVER)

export interface Connection {
    host: string,
    sshKeyId?: string,
    sshPort?: number,
    keeperPort?: number,
    dbPort?: number,
}

export interface NodeResponse {
    connection: Connection,
    keeper: KeeperResponse,
}

export interface NodeRequest {
    connection: Connection,
    vaultId?: string,
    certs?: Certs,
    body?: any,
}

export interface SshRequest {
    connection: Connection,
}

export interface DockerRequest {
    connection: Connection,
    image?: string,
    container?: string,
    options?: string,
}

export interface DockerLogsRequest {
    connection: Connection,
    container: string,
    tail?: number,
}

export interface DockerResult {
    stdout: string,
    stderr: string,
    exitCode: number,
}

export enum KeeperStatus {
    Active = "ACTIVE",
    Paused = "PAUSED",
}

export interface Keeper {
    host: string,
    port: number,
    name?: string,
    status?: KeeperStatus,
}

export type Role = "leader" | "replica" | "unknown";

export interface KeeperResponse {
    name?: string,
    status?: KeeperStatus,
    state: string,
    role: Role,
    lag: number,
    pendingRestart: boolean,
    scheduledSwitchover?: ScheduledSwitchover,
    scheduledRestart?: ScheduledRestart,
    tags?: {[key: string]: any},
    discoveredHost: string,
    discoveredKeeperPort: number,
    discoveredDbPort: number,
}

export interface ScheduledSwitchover {
    at: string,
    to: string,
}

export interface ScheduledRestart {
    at: string,
    pendingRestart: boolean,
}

// SPECIFIC (WEB)

export enum NodeTabType {QUERY, CHART}
export interface NodeTab {
    label: string,
    body: (connection: QueryConnection) => ReactNode,
    info?: ReactNode,
}
