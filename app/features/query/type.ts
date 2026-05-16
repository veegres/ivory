import {Certs} from "../cert/type"
import {Config} from "../database/type"

// COMMON (WEB AND SERVER)

export enum Type {
    BLOAT,
    ACTIVITY,
    REPLICATION,
    STATISTIC,
    OTHER,
    CONSOLE,
}

export enum VarietyType {
    DatabaseSensitive,
    MasterOnly,
    ReplicaRecommended,
}

export enum CreationType {
    Manual = "manual",
    System = "system",
}

export enum ChartType {
    Databases = "Databases",
    Connections = "Connections",
    DatabaseSize = "Database Size",
    DatabaseUptime = "Database Uptime",
    Schemas = "Schemas",
    TablesSize = "Tables Size",
    IndexesSize = "Indexes Size",
    TotalSize = "Total Size",
}

export interface Request {
    name: string,
    type?: Type,
    description?: string,
    query: string,
    varieties?: VarietyType[],
    params?: string[],
}

export interface Response {
    id: string,
    name: string,
    type: Type,
    creation: CreationType,
    varieties?: VarietyType[],
    params?: string[],
    description?: string,
    default: string,
    custom: string,
    createdAt: number,
}

export interface Connection {
    db: Config,
    certs?: Certs,
    vaultId?: string,
}

export interface TemplateRequest {
    connection: Connection,
    queryUuid?: string,
    options?: DbOptions,
}

export interface ConsoleRequest {
    connection: Connection,
    query: string,
    options?: DbOptions,
}

export interface KillRequest {
    connection: Connection,
    pid: number,
}

export interface ChartRequest {
    connection: Connection,
    type: ChartType,
}


export interface DatabasesRequest {
    connection: Connection,
    name: string,
}

export interface SchemasRequest {
    connection: Connection,
    name: string,
}

export interface TablesRequest {
    connection: Connection,
    schema: string,
    name: string,
}

export interface Chart {
    name: string,
    value: any,
}

export interface DbOptions {
    params?: any[],
    trim?: boolean,
    limit?: string,
}

export interface DbRow {
    name: string,
    dataType: string,
    dataTypeOID: number,
}

export interface DbResponse {
    fields: DbRow[],
    rows: any[][],
    url: string,
    startTime: number,
    endTime: number,
    options?: DbOptions,
}

// SPECIFIC (WEB)

export interface RunRequest {
    connection: Connection,
    queryUuid?: string,
    query: string,
    options?: DbOptions,
}
