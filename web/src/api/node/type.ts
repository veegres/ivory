import {ReactNode} from "react"

import {Certs} from "../cert/type"
import {KeeperResponse as BaseKeeperResponse, Plugin as KeeperPlugin, Status as KeeperStatus} from "../keeper/type"
import {Docker, Metrics} from "../os/type"
import {Connection as QueryConnection} from "../query/type"

// COMMON (WEB AND SERVER)

export {KeeperPlugin,KeeperStatus}

export interface KeeperConnection {
    host: string,
    port: number,
}

export interface KeeperRequest extends KeeperConnection {
    plugin: KeeperPlugin,
    vaultId?: string,
    certs?: Certs,
    body?: any,
}

export type KeeperResponse = BaseKeeperResponse

export interface SshConnection {
    host: string,
    port: number,
    vaultId: string,
}

export type MetricsResponse = Metrics

export interface DockerRequest {
    connection: SshConnection,
    image?: string,
    container?: string,
    options?: string,
}

export interface DockerLogsRequest {
    connection: SshConnection,
    container: string,
    tail?: number,
}

export type DockerResponse = Docker

// SPECIFIC (WEB)

export enum NodeTabType {QUERY, MONITOR}
export interface NodeTab {
    label: string,
    body: (queryCon?: QueryConnection, sshCon?: SshConnection) => ReactNode,
    info?: ReactNode,
}
