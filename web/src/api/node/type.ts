import {ReactNode} from "react"

import {Certs} from "../cert/type"
import {Container, Metrics} from "../cloud/type"
import {KeeperResponse as BaseKeeperResponse, Plugin as KeeperPlugin, Status as KeeperStatus} from "../keeper/type"
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

export interface CloudConnection {
    host: string,
    port: number,
    vaultId: string,
}

export type MetricsResponse = Metrics

export interface ContainerRequest {
    connection: CloudConnection,
    image?: string,
    container?: string,
    options?: string,
}

export interface ContainerLogsRequest {
    connection: CloudConnection,
    container: string,
    tail?: number,
}

export type ContainerResponse = Container

// SPECIFIC (WEB)

export enum NodeTabType {QUERY, MONITOR}
export interface NodeTab {
    label: string,
    body: (queryCon?: QueryConnection, cloudCon?: CloudConnection) => ReactNode,
    info?: ReactNode,
}
