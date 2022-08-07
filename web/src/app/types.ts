import {CSSProperties, ReactNode} from "react";

export interface Response<TData, TError = {}> {
    response: TData
    error: TError
}

// NODE
export interface InstanceResponse {
    name: string,
    timeline?: number,
    lag: number,
    state: string,
    host: string,
    role: string,
    port: number,
    api_url: string
}

export interface Instance extends InstanceResponse{
    api_domain: string
    isLeader: boolean
    credId?: string
}

export interface InstanceOverview {
    database_system_identifier: string,
    pending_restart?: boolean
    postmaster_start_time?: string,
    timeline?: number,
    cluster_unlocked?: boolean,
    patroni: { scope: string, version: string },
    state: string,
    role: string,
    xlog?: {
        received_location: number,
        replayed_timestamp: string,
        paused: boolean,
        replayed_location: number
    }
    server_version?: string
}

// CLUSTER
export interface Cluster {
    name: string,
    certsId?: string,
    patroniCredId?: string,
    postgresCredId?: string,
    nodes: string[],
    tags?: string[]
}

export interface ClusterMap {
    [name: string]: string[]
}

export interface ClusterTabs {
    [key: number]: { body: ReactNode, info?: ReactNode }
}

// SECRET
export interface SecretStatus {
    key: boolean
    ref: boolean
}

export interface SecretSetRequest {
    key: string
    ref: string
}

export interface SecretUpdateRequest {
    previousKey: string
    newKey: string
}

// CREDENTIAL
export interface Credential {
    username: string
    password: string
    type: CredentialType
}

export enum CredentialType {
    POSTGRES,
    PATRONI
}

export interface CredentialMap {
    [uuid: string]: Credential
}

// BLOAT
export interface Connection {
    host: string
    port: number
    credId: string
}

export interface Target {
    dbName?: string
    schema?: string
    table?: string
    excludeSchema?: string
    excludeTable?: string
}

export interface CompactTableRequest {
    cluster: string
    connection: Connection
    target?: Target
    ratio?: number
}

export interface CompactTable {
    uuid: string,
    status: JobStatus,
    credentialId: string,
    command: string,
    commandArgs: string,
    logsPath: string
}

export enum JobStatus {
    PENDING,
    RUNNING,
    FINISHED,
    FAILED,
    STOPPED,
    UNKNOWN
}

export enum EventStream {
    START = "start",
    END = "end"
}

// COMMON
export interface Style {
    [key: string]: CSSProperties
}

export interface ColorsMap {
    [name: string]: 'success' | 'primary'
}
