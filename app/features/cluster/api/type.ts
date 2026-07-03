import {ReactNode} from "react"

import {Certs} from "../../cert/api/type"
import {Feature} from "../../feature"
import {KeeperOneResponse, KeeperPlugin} from "../../node/api/type"
import {DbPlugin} from "../../query/api/type"

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

export interface ImageConfig {
    cluster: string,
    dcs: string,
    dbPass: string,
    dbUser: string,
}

export interface DeployCommonConfig extends ImageConfig {
    sshPass: string,
    sshUser: string,
}

export interface DeployRequest {
    uri: string,
    parallel: boolean
    nodeRawImageOptions: {[key: string]: string},
    nodeConfig: NodeConfig[],
    commonConfig: DeployCommonConfig,
    clusterOptions: Options
}

export interface NodeOverview {
    [domain: string]: Node,
}

export interface Node {
    config: NodeConfig,
    keeper: KeeperOneResponse,
    warnings: string[],
}

// SPECIFIC (WEB)

export interface InterpolatedOptions extends ImageConfig, NodeConfig {}

export interface ClusterTab {
    label: string,
    body: (cluster: Cluster, mainNode?: Node, nodes?: NodeOverview) => ReactNode,
    feature: Feature,
    info?: ReactNode,
}
