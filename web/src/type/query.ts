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

export interface Query {
    id: string,
    name: string,
    type: QueryType,
    creation: QueryCreation,
    description: string,
    default: string,
    custom: string,
}

export interface QueryRequest {
    name?: string,
    type?: QueryType,
    description?: string,
    query: string,
}

export interface QueryPostgresRequest {
    credentialId?: string,
    db: Database,
}

export interface QueryRunRequest extends QueryPostgresRequest {
    queryUuid: string,
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
}

// SPECIFIC (WEB)
