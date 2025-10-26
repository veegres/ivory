import {Certs} from "../cert/type"

// COMMON (WEB AND SERVER)

export enum QueryCreation {
    Manual = "manual",
    System = "system",
}

export interface Query {
    id: string,
    name: string,
    type: QueryType,
    creation: QueryCreation,
    varieties?: QueryVariety[],
    params?: string[],
    description?: string,
    default: string,
    custom: string,
    createAt: number,
}

export interface QueryConnection {
    db: Database,
    certs?: Certs,
    credentialId?: string,
}

export interface QueryRunRequest {
    connection: QueryConnection,
    queryUuid?: string,
    query?: string,
    queryOptions?: QueryOptions,
}

export interface QueryKillRequest {
    connection: QueryConnection,
    pid: number,
}

export interface QueryChartRequest {
    connection: QueryConnection,
    type: QueryChartType,
}


export interface QueryDatabasesRequest {
    connection: QueryConnection,
    name: string,
}

export interface QuerySchemasRequest {
    connection: QueryConnection,
    name: string,
}

export interface QueryTablesRequest {
    connection: QueryConnection,
    schema: string,
    name: string,
}

// POSTGRES CLIENT (WEB AND SERVER)

export interface Database {
    host: string,
    port: number,
    name?: string,
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

export interface QueryChart {
    name: string,
    value: any,
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


export interface QueryRequest {
    type: QueryType,
    query: string,
    name: string,
    description?: string,
    varieties?: QueryVariety[],
    params?: string[],
}

// SPECIFIC (WEB)
