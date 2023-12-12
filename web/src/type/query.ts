import {Database} from "./common";

// COMMON (WEB AND SERVER)

export enum QueryType {
    BLOAT,
    ACTIVITY,
    REPLICATION,
    STATISTIC,
    OTHER,
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
    description: string,
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

export interface QueryPostgresRequest {
    credentialId?: string,
    db: Database,
}

export interface QueryRunRequest extends QueryPostgresRequest {
    queryUuid: string,
    queryParams?: string[],
}

export interface QueryKillRequest extends QueryPostgresRequest {
    pid: number,
}

export interface QueryChartRequest extends QueryPostgresRequest {}


export interface QueryDatabasesRequest extends QueryPostgresRequest {
    name: string,
}

export interface QuerySchemasRequest extends QueryPostgresRequest {
    name: string,
}

export interface QueryTablesRequest extends QueryPostgresRequest {
    schema: string,
    name: string,
}

export interface QueryField {
    name: string,
    dataType: string,
    dataTypeOID: string,
}

export interface QueryChart {
    name: string,
    value: any,
}

export interface QueryRunResponse {
    fields: QueryField[],
    rows: any[][],
    url: string,
}

// SPECIFIC (WEB)
