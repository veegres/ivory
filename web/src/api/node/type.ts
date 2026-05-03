import {ReactNode} from "react"

import {Certs} from "../cert/type"
import {KeeperResponse as BaseKeeperResponse, Plugin as KeeperPlugin, Status as KeeperStatus} from "../keeper/type"
import {Connection as QueryConnection} from "../query/type"

// COMMON (WEB AND SERVER)

export {KeeperPlugin,KeeperStatus}

export interface Connection {
    host: string,
    sshPort?: number,
    keeperPort?: number,
    dbPort?: number,
}

export interface KeeperConnection {
    host: string,
    port: number,
    plugin: KeeperPlugin,
    vaultId?: string,
    certs?: Certs,
    body?: any,
}

export interface KeeperResponse {
    connection: Connection,
    keeper: BaseKeeperResponse,
}

export interface SshConnection {
    host: string,
    port: number,
    vaultId: string,
}

export interface CpuMetrics {
    totalTicks: number,
    idleTicks: number,
}

export interface MemoryMetrics {
    totalBytes: number,
    availableBytes: number,
}

export interface NetworkMetrics {
    receivedBytes: number,
    transmittedBytes: number,
}

export interface NodeMetrics {
    cpu: CpuMetrics,
    memory: MemoryMetrics,
    network: NetworkMetrics,
}

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

export interface DockerResult {
    stdout: string,
    stderr: string,
    exitCode: number,
}

// SPECIFIC (WEB)

export enum NodeTabType {QUERY, MONITOR}
export interface NodeTab {
    label: string,
    body: (queryCon?: QueryConnection, sshCon?: SshConnection) => ReactNode,
    info?: ReactNode,
}
