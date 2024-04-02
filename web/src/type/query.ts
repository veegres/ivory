import {Database} from "./general";
import {Certs} from "./cluster";

// COMMON (WEB AND SERVER)

export enum QueryType {
    BLOAT,
    ACTIVITY,
    REPLICATION,
    STATISTIC,
    OTHER,
    CONSOLE,
}

export enum QueryCreation {
    Manual = "manual",
    System = "system",
}

export enum QueryVariety {
    DatabaseSensitive,
    MasterOnly,
    ReplicaRecommended,
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

export interface QueryRequest {
    type: QueryType,
    query: string,
    name: string,
    description?: string,
    varieties?: QueryVariety[],
    params?: string[],
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
    queryParams?: string[],
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

export interface QueryField {
    name: string,
    dataType: string,
    dataTypeOID: number,
}

export enum QueryChartType {
    Databases = "Databases",
    Connections = "Connections",
    DatabaseSize = "DatabaseSize",
    DatabaseUptime = "DatabaseUptime",
    Schemas = "Schemas",
    TablesSize = "TablesSize",
    IndexesSize = "IndexesSize",
    TotalSize = "TotalSize",
}

export interface QueryChart {
    name: string,
    value: any,
}

export interface QueryFields {
    fields: QueryField[],
    rows: any[][],
    url: string,
    startTime: number,
    endTime: number,
}

// SPECIFIC (WEB)
