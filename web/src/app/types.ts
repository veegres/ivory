import {CSSProperties} from "react";

// SERVER
export interface Node {
    name: string,
    timeline: number,
    lag: number,
    state: string,
    host: string,
    role: string,
    port: number,
    api_url: string
}

export interface NodeOverview {
    database_system_identifier: string,
    postmaster_start_time: string,
    timeline: number,
    cluster_unlocked: boolean,
    patroni: { scope: string, version: string },
    state: string,
    role: string,
    xlog: {
        received_location: number,
        replayed_timestamp: string,
        paused: boolean,
        replayed_location: number
    }
    server_version: string
}

export interface Cluster {
    name: string,
    nodes: string[]
}

export interface ClusterMap {
    [name: string]: string[]
}

export interface Auth {
    username: string
    password: string
}

export interface Connection extends Auth {
    host: string
    port: number
}

export interface Target {
    dbName?: string
    schema?: string
    table?: string
    excludeSchema?: string
    excludeTable?: string
}

export interface CompactTableRequest {
    connection: Connection
    target?: Target
    ratio?: number
}

export interface Response<TData, TError = {}> {
    response: TData
    error: TError
}

export interface CompactTable {
    uuid: string,
    status: JobStatus,
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

// COMMON
export interface Style {
    [key: string]: CSSProperties
}

export enum EventStream {
    START = "start",
    END = "end"
}
