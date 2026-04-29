import {ReactNode} from "react"

import {Certs} from "../cert/type"
import {KeeperResponse as BaseKeeperResponse, Plugin as KeeperPlugin, Status as KeeperStatus} from "../keeper/type"
import {Connection as QueryConnection} from "../query/type"

// COMMON (WEB AND SERVER)

export {KeeperPlugin,KeeperStatus}

export interface Connection {
    host: string,
    sshKeyId?: string,
    sshPort?: number,
    keeperPort?: number,
    dbPort?: number,
}

export interface KeeperResponse {
    connection: Connection,
    keeper: BaseKeeperResponse,
}

export interface KeeperRequest {
    host: string,
    port: number,
    plugin: KeeperPlugin,
    vaultId?: string,
    certs?: Certs,
    body?: any,
}

export interface SshRequest {
    connection: Connection,
}

export interface DockerRequest {
    connection: Connection,
    image?: string,
    container?: string,
    options?: string,
}

export interface DockerLogsRequest {
    connection: Connection,
    container: string,
    tail?: number,
}

export interface DockerResult {
    stdout: string,
    stderr: string,
    exitCode: number,
}

// SPECIFIC (WEB)

export enum NodeTabType {QUERY, MONITOR}
export interface NodeTab {
    label: string,
    body: (queryCon: QueryConnection, sshCon: Connection) => ReactNode,
    info?: ReactNode,
}
