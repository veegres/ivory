// COMMON (WEB AND SERVER)

import {Certs} from "./cert/type"

export interface Database {
    host: string,
    port: number,
    name?: string,
    schema?: string,
}

export interface ConnectionRequest {
    db: Database,
    certs?: Certs,
    credentialId?: string,
}


export interface QueryOptions {
    params?: string[],
    trim?: boolean,
    limit?: string,
}

export interface QueryField {
    name: string,
    dataType: string,
    dataTypeOID: number,
}

export interface QueryFields {
    fields: QueryField[],
    rows: any[][],
    url: string,
    startTime: number,
    endTime: number,
    options?: QueryOptions,
}

export enum QueryChartType {
    Databases = "Databases",
    Connections = "Connections",
    DatabaseSize = "Database Size",
    DatabaseUptime = "Database Uptime",
    Schemas = "Schemas",
    TablesSize = "Tables Size",
    IndexesSize = "Indexes Size",
    TotalSize = "Total Size",
}

export enum QueryType {
    BLOAT,
    ACTIVITY,
    REPLICATION,
    STATISTIC,
    OTHER,
    CONSOLE,
}

export enum QueryVariety {
    DatabaseSensitive,
    MasterOnly,
    ReplicaRecommended,
}


export interface Query {
    type: QueryType,
    query: string,
    name: string,
    description?: string,
    varieties?: QueryVariety[],
    params?: string[],
}
