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

export interface Connection {
    host: string
    port: number
    username: string
    password: string
}

export interface Target {
    dbName?: string
    schema?: string
    table?: string
    excludeSchema?: string
    excludeTable?: string
}

export interface PgCompactTable {
    connection: Connection
    target?: Target
    ratio?: number
}

export interface GoResponse<TData, TError = {}> {
    response: TData
    error: TError
}

// COMMON
export interface Style {
    [key: string]: CSSProperties
}
