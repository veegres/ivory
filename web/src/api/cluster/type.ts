import {ReactNode} from "react"

import {Certs} from "../cert/type"
import {Plugin as DbPlugin} from "../database/type"
import {Feature} from "../feature"
import {Plugin as KeeperPlugin} from "../keeper/type"
import {KeeperResponse} from "../node/type"

// COMMON (WEB AND SERVER)

export interface Plugins {
    keeper: KeeperPlugin,
    database: DbPlugin,
}

export interface Options {
    plugins: Plugins,
    tls: Tls,
    certs: Certs,
    vaults: Vaults,
    tags: string[],
}

export interface NodeConfig {
    host: string,
    sshPort?: number,
    keeperPort?: number,
    dbPort?: number,
}

export interface Cluster extends Options {
    name: string,
    nodes: NodeConfig[],
    nodesOverview?: NodeOverview,
}

export interface AutoRequest extends Options {
    name: string,
    host: string,
    port: number,
}

export interface Tls {
    keeper: boolean,
    database: boolean,
}

export interface Vaults {
    keeperId?: string,
    databaseId?: string,
    sshKeyId?: string,
}

export interface Overview {
    nodes: NodeOverview,
    features: Feature[],
}

export interface NodeOverview {
    [domain: string]: Node,
}

export interface Node {
    config: NodeConfig,
    keeper: KeeperResponse,
    warnings: string[],
}

// SPECIFIC (WEB)

export interface ClusterTab {
    label: string,
    body: (cluster: Cluster, mainNode?: Node, nodes?: NodeOverview) => ReactNode,
    feature: Feature,
    info?: ReactNode,
}
