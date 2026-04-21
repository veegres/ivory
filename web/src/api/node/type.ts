
import {ReactNode} from "react"

import {Certs} from "../cert/type"
import {ConnectionRequest} from "../postgres"

// COMMON (WEB AND SERVER)

export interface NodeRequest {
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
