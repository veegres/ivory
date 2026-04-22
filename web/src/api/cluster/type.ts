import {ReactNode} from "react"

import {Certs} from "../cert/type"
import {
    Keeper,
    NodeScheduledRestart,
    NodeScheduledSwitchover, Role,
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
    keepers: Keeper[],
    keepersOverview: NodeOverview,
}

export interface ClusterAuto extends ClusterOptions {
    name: string,
    node: Keeper,
}

export interface ClusterTls {
    keeper: boolean,
    database: boolean,
}

export interface Credentials {
    patroniId?: string,
    postgresId?: string,
}


export interface ClusterOverview {
    nodes: NodeOverview,
    detectedBy?: Keeper,
    mainNode?: Node,
}

export interface NodeOverview {
    [keeper: string]: Node | undefined
}

export interface Node {
    state: string,
    role: Role,
    lag: number,
    pendingRestart: boolean,
    database: Database,
    keeper: Keeper,
    scheduledSwitchover?: NodeScheduledSwitchover,
    scheduledRestart?: NodeScheduledRestart,
    tags?: {[key: string]: any},
    inCluster: boolean,
    inKeeper: boolean,
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
