import {ReactNode} from "react"

import {Certs} from "../../cert/api/CertType"
import {DeployCredentials} from "../../deployment/api/DeploymentType"
import {Feature} from "../../Feature"
import {DeploySteps, KeeperLag, KeeperOneRequest, KeeperOneResponse, KeeperPlugin, NodeCheckStep, PlatformPlugin, PlatformVaultConnection} from "../../node/api/NodeType"
import {Connection as QueryConnection, DbPlugin, QueryCheckStep} from "../../query/api/QueryType"

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
    // name is the node's own name, unique within the cluster and independent of
    // its host: it is the deployment's identity ({{name}}, --name) and the name
    // the platform addresses the container by, so it is always set - a node the
    // keeper reports no name for falls back to its host
    name: string,
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
    lag: KeeperLag,
}

// NOTE: a half is missing when the node has no port or vault to reach it with
export interface NodeConnections {
    keeper?: KeeperOneRequest,
    platform?: PlatformVaultConnection,
    database?: QueryConnection,
}

export interface ClusterConnections {
    [node: string]: NodeConnections,
}

// NOTE: a list is missing where NodeConnections has no connection for it
export interface ClusterNodeCheck {
    keeper?: NodeCheckStep[],
    platform?: NodeCheckStep[],
    database?: QueryCheckStep[],
}

export interface ClusterCheck {
    [node: string]: ClusterNodeCheck,
}

export interface DeployCommonConfig {
    cluster: string,
    dcs: string,
    sshUser: string,
    sshPass: string,
    keeperUser: string,
    keeperPass: string,
    dbUser: string,
    dbPass: string,
}

// DeployPreviewCredentials is the pair of previews a node card needs - what
// may be shown of the keeper's credentials and of the database's.
export interface DeployPreviewCredentials {
    keeper: DeployCredentials,
    database: DeployCredentials,
}

// DeployNode pairs one node with the command that deploys it, for the length
// of one request only - the command is never persisted on the cluster.
export interface DeployNode extends NodeConfig {
    command: string,
    postScripts?: string[],
}

export interface DeployRequest {
    parallel: boolean,
    nodes: DeployNode[],
    platform: PlatformPlugin,
    commonConfig: DeployCommonConfig,
    clusterOptions: Options
}

export interface DeployResponse {
    logs: string[],
    nodes: DeployNodeSteps,
}

export interface DeployNodeSteps {
    [name: string]: DeploySteps,
}

export interface NodeOverview {
    [domain: string]: Node,
}

export interface Node {
    config: NodeConfig,
    keeper: KeeperOneResponse,
    // NOTE: Ivory's own, the engine's are in keeper.warnings
    warnings: string[],
}

// SPECIFIC (WEB)

export interface ClusterTab {
    label: string,
    body: (cluster: Cluster, mainNode?: Node, nodes?: NodeOverview) => ReactNode,
    feature: Feature,
    info?: ReactNode,
}
