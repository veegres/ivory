
import {ReactNode} from "react"

import {Certs} from "../cert/type"
import {ConnectionRequest} from "../postgres"

// COMMON (WEB AND SERVER)

export interface VM {
    id: string,
    name: string,
    host: string,
    sshPort: number,
    username: string,
    sshKey: string,
}

export interface NodeConnection {
    vmId?: string,
    host: string,
    sshPort?: number,
    keeperPort: number,
    dbPort?: number,
}

export interface NodeRequest {
    connection: NodeConnection,
    credentialId?: string,
    certs?: Certs,
    body?: any,
}

export interface NodeAutoRequest {
    connections: NodeConnection[],
    credentialId?: string,
    certs?: Certs,
    body?: any,
}

export interface Node {
    connection: NodeConnection,
    response: KeeperResponse,
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
    scheduledSwitchover?: NodeScheduledSwitchover,
    scheduledRestart?: NodeScheduledRestart,
    tags?: {[key: string]: any},
    discoveredHost: string,
    discoveredKeeperPort: number,
    discoveredDbPort: number,
}

export interface NodeScheduledSwitchover {
    at: string,
    to: string,
}

export interface NodeScheduledRestart {
    at: string,
    pendingRestart: boolean,
}

// SPECIFIC (WEB)

export enum NodeTabType {QUERY, CHART}
export interface NodeTab {
    label: string,
    body: (connection: ConnectionRequest) => ReactNode,
    info?: ReactNode,
}
