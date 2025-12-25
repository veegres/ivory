import {ConnectionRequest, QueryChartType, QueryOptions, QueryType, QueryVariety} from "../postgres"

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

export interface QueryTemplateRequest {
    connection: ConnectionRequest,
    queryUuid?: string,
    queryOptions?: QueryOptions,
}

export interface QueryConsoleRequest {
    connection: ConnectionRequest,
    query: string,
    queryOptions?: QueryOptions,
}

export interface QueryKillRequest {
    connection: ConnectionRequest,
    pid: number,
}

export interface QueryChartRequest {
    connection: ConnectionRequest,
    type: QueryChartType,
}


export interface QueryDatabasesRequest {
    connection: ConnectionRequest,
    name: string,
}

export interface QuerySchemasRequest {
    connection: ConnectionRequest,
    name: string,
}

export interface QueryTablesRequest {
    connection: ConnectionRequest,
    schema: string,
    name: string,
}

export interface QueryChart {
    name: string,
    value: any,
}

// SPECIFIC (WEB)

export interface QueryRunRequest {
    connection: ConnectionRequest,
    queryUuid?: string,
    query: string,
    queryOptions?: QueryOptions,
}