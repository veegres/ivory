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

export interface DeployCommonConfig {
    cluster: string,
    sshUser: string,
    sshPass: string,
    dbUser: string,
    dbPass: string,
}

export interface DeployNode extends NodeConfig {
    // options overrides the rendered options template for this node
    options?: string,
}

// DeployRequest describes a deployment intent: node ports, the image, aux
// ports, the DCS value and the per-node options are resolved server-side
// from the keeper plugin's spec unless explicitly provided.
export interface DeployRequest {
    parallel: boolean,
    singleHost: boolean,
    image?: string,
    nodes: DeployNode[],
    values: {[key: string]: string},
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

export interface ClusterTab {
    label: string,
    body: (cluster: Cluster, mainNode?: Node, nodes?: NodeOverview) => ReactNode,
    feature: Feature,
    info?: ReactNode,
}
