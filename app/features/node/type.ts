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

export interface KeeperRequest extends KeeperConnection {
    plugin: KeeperPlugin,
    vaultId?: string,
    certs?: Certs,
    body?: any,
}

export type KeeperResponse = BaseKeeperResponse

export interface PlatformConnection {
    host: string,
    port: number,
    vaultId: string,
}

export type MetricsResponse = Metrics

export interface PlatformDeployRequest {
    connection: PlatformConnection,
    image?: string,
    name?: string,
    options?: string,
}

export interface PlatformLogsRequest {
    connection: PlatformConnection,
    name: string,
    tail?: number,
}

// SPECIFIC (WEB)

export enum NodeTabType {QUERY, MONITOR}
export interface NodeTab {
    label: string,
    body: (queryCon?: QueryConnection, platformCon?: PlatformConnection) => ReactNode,
    info?: ReactNode,
}
