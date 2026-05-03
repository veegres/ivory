import {ReactNode} from "react"

import {Certs} from "../cert/type"
import {Plugin as DbPlugin} from "../database/type"
import {Feature} from "../feature"
import {Plugin as KeeperPlugin} from "../keeper/type"
import {Connection, KeeperResponse} from "../node/type"

// COMMON (WEB AND SERVER)

export interface Plugins {
    keeper: KeeperPlugin,
    database: DbPlugin,
}

export interface ClusterOptions {
    plugins: Plugins,
    tls: ClusterTls,
    certs: Certs,
    vaults: Vault,
    tags: string[],
}

export interface Cluster extends ClusterOptions {
    name: string,
    nodes: Connection[],
    nodesOverview?: NodeOverview,
}

export interface ClusterAuto extends ClusterOptions {
    name: string,
    host: string,
    port: number,
}

export interface ClusterTls {
    keeper: boolean,
    database: boolean,
}

export interface Vault {
    keeperId?: string,
    databaseId?: string,
    sshKeyId?: string,
}

export interface ClusterOverview {
    nodes: NodeOverview,
    detectedDomain: string,
    features: Feature[],
}

export interface NodeOverview {
    [domain: string]: Node,
}

export interface Node extends KeeperResponse {
    warnings: string[],
}

// SPECIFIC (WEB)

export interface ActiveCluster {
    cluster: Cluster,
    warning: boolean,
    manualKeeper?: string,
}

export interface ClusterTab {
    label: string,
    body: (cluster: Cluster, mainNode?: Node, nodes?: NodeOverview) => ReactNode,
    feature: Feature,
    info?: ReactNode,
}
