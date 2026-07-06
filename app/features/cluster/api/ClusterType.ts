import {ReactNode} from "react"

import {Certs} from "../../cert/api/CertType"
import {Feature} from "../../Feature"
import {KeeperOneResponse, KeeperPlugin} from "../../node/api/NodeType"
import {DbPlugin} from "../../query/api/QueryType"

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
    // NOTE: a feature absent from the map is not a plugin capability at all
    // and should be treated as unrestricted, only `false` means unsupported
    features: Partial<Record<Feature, boolean>>,
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
