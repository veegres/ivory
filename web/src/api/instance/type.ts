
import {ReactNode} from "react"

import {Certs} from "../cert/type"
import {ConnectionRequest} from "../postgres"

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

export interface InstanceScheduledSwitchover {
    at: string,
    to: string,
}

export interface InstanceScheduledRestart {
    at: string,
    pendingRestart: boolean,
}

// SPECIFIC (WEB)

export enum InstanceTabType {QUERY, CHART}
export interface InstanceTab {
    label: string,
    body: (connection: ConnectionRequest) => ReactNode,
    info?: ReactNode,
}
