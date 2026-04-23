import {ReactNode} from "react"

import {Certs} from "../cert/type"
import {
    Connection,
    Keeper,
    NodeResponse as BaseNode,
} from "../node/type"
import {Permission} from "../permission/type"

// COMMON (WEB AND SERVER)

export interface ClusterOptions {
    tls: ClusterTls,
    certs: Certs,
    vaults: Vault,
    tags?: string[],
}

export interface Cluster extends ClusterOptions {
    name: string,
    nodes: Connection[],
    nodesOverview: NodeOverview,
}

export interface ClusterAuto extends ClusterOptions {
    name: string,
    node: Connection,
}

export interface ClusterTls {
    keeper: boolean,
    database: boolean,
}

export interface Vault {
    keeperId?: string,
    databaseId?: string,
}


export interface ClusterOverview {
    nodes: NodeOverview,
    detectedBy?: Keeper,
    mainNode?: Node,
}

export interface NodeOverview {
    [domain: string]: Node | undefined
}

export interface Node extends BaseNode {
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
