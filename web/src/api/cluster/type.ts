import {ReactNode} from "react"

import {Certs} from "../cert/type"
import {
    NodeScheduledRestart,
    NodeScheduledSwitchover, Role,
    Sidecar,
} from "../node/type"
import {Permission} from "../permission/type"
import {Database} from "../postgres"

// COMMON (WEB AND SERVER)

export interface ClusterOptions {
    tls: ClusterTls,
    certs: Certs,
    credentials: Credentials,
    tags?: string[],
}

export interface Cluster extends ClusterOptions {
    name: string,
    sidecars: Sidecar[],
    sidecarsOverview: NodeOverview,
}

export interface ClusterAuto extends ClusterOptions {
    name: string,
    node: Sidecar,
}

export interface ClusterTls {
    sidecar: boolean,
    database: boolean,
}

export interface Credentials {
    patroniId?: string,
    postgresId?: string,
}


export interface ClusterOverview {
    nodes: NodeOverview,
    detectedBy?: Sidecar,
    mainNode?: Node,
}

export interface NodeOverview {
    [sidecar: string]: Node | undefined
}

export interface Node {
    state: string,
    role: Role,
    lag: number,
    pendingRestart: boolean,
    database: Database,
    sidecar: Sidecar,
    scheduledSwitchover?: NodeScheduledSwitchover,
    scheduledRestart?: NodeScheduledRestart,
    tags?: {[key: string]: any},
    inCluster: boolean,
    inSidecar: boolean,
}

// SPECIFIC (WEB)

export interface ActiveCluster {
    cluster: Cluster,
    detectBy?: Node,
    warning: boolean,
}

export interface ActiveNode {
    [cluster: string]: string | undefined
}

export interface ClusterTab {
    label: string,
    body: (cluster: Cluster, overview?: ClusterOverview) => ReactNode,
    permission: Permission,
    info?: ReactNode,
}
