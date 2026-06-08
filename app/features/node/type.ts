import {ReactNode} from "react"

import {Certs} from "../cert/type"
import {KeeperResponse as BaseKeeperResponse, Plugin as KeeperPlugin, Status as KeeperStatus} from "../keeper/type"
import {Metrics} from "../platform/type"
import {Connection as QueryConnection} from "../query/type"

// COMMON (WEB AND SERVER)

export {KeeperPlugin,KeeperStatus}

export interface KeeperConnection {
    host: string,
    port: number,
}

export interface KeeperOneRequest extends KeeperConnection {
    plugin: KeeperPlugin,
    vaultId?: string,
    certs?: Certs,
    body?: any,
}

export type KeeperOneResponse = BaseKeeperResponse

export interface KeeperMultiRequest {
    connections: KeeperConnection[],
    body?: any,
}

export interface KeeperMultiResponse {
    connection: KeeperConnection,
    response: BaseKeeperResponse[],
    error?: string,
}

export interface PlatformConnection {
    host: string,
    port: number,
    vaultId: string,
}

export interface PlatformCredConnection {
    host: string,
    port: number,
    username: string,
    password?: string,
}

export type PlatformMetricsResponse = Metrics

export interface PlatformMetricsRequest {
    host: string,
    port: number,
    vaultId: string,
}

export interface PlatformCopyIdRequest {
    host: string,
    port: number,
    username: string,
    password?: string,
    publicKey: string,
}

export interface PlatformUpRequest {
    connection: PlatformConnection,
    image: string,
    name: string,
    options?: string,
}

export interface PlatformLogsRequest {
    connection: PlatformConnection,
    name: string,
    tail?: number,
    follow?: boolean,
}

export interface PlatformActionRequest {
    connection: PlatformConnection,
    name: string,
}

// SPECIFIC (WEB)

export enum NodeTabType {QUERY, MONITOR}
export interface NodeTab {
    label: string,
    body: (queryCon?: QueryConnection, platformCon?: PlatformConnection) => ReactNode,
    info?: ReactNode,
}
