
import {ReactNode} from "react";
import {Database, QueryConnection} from "../query/type";
import {Certs} from "../cert/type";

// COMMON (WEB AND SERVER)

export interface InstanceRequest {
    sidecar: Sidecar,
    credentialId?: string,
    certs?: Certs,
    body?: any,
}

// PATRONI CLIENT (WEB AND SERVER)

export enum SidecarStatus {
    Active = "ACTIVE",
    Paused = "PAUSED",
}

export interface Sidecar {
    host: string,
    port: number,
    name?: string,
    status?: SidecarStatus,
}

export type Role = "leader" | "replica" | "unknown";
export interface Instance {
    state: string,
    role: Role,
    lag: number,
    pendingRestart: boolean,
    database: Database,
    sidecar: Sidecar,
    scheduledSwitchover?: InstanceScheduledSwitchover,
    scheduledRestart?: InstanceScheduledRestart,
    tags?: {[key: string]: any},
}

export interface InstanceScheduledSwitchover {
    at: string,
    to: string,
}

export interface InstanceScheduledRestart {
    at: string,
    pendingRestart: boolean,
}

// SPECIFIC (WEB)

// TODO should we return map from the server?
export interface InstanceMap {
    [instance: string]: InstanceWeb,
}

export enum InstanceTabType {QUERY, CHART}
export interface InstanceTab {
    label: string,
    body: (connection: QueryConnection) => ReactNode,
    info?: ReactNode,
}

export interface InstanceWeb extends Instance {
    leader: boolean,
    inCluster: boolean,
    inInstances: boolean,
}
