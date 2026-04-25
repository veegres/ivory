import {Certs} from "../cert/type"
import {Database, QueryChartType, QueryOptions, QueryType, QueryVariety} from "../database/type"

// COMMON (WEB AND SERVER)

export {QueryType}

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

export interface Connection {
    db: Database,
    certs?: Certs,
    vaultId?: string,
}

export interface QueryTemplateRequest {
    connection: Connection,
    queryUuid?: string,
    queryOptions?: QueryOptions,
}

export interface QueryConsoleRequest {
    connection: Connection,
    query: string,
    queryOptions?: QueryOptions,
}

export interface QueryKillRequest {
    connection: Connection,
    pid: number,
}

export interface QueryChartRequest {
    connection: Connection,
    type: QueryChartType,
}


export interface QueryDatabasesRequest {
    connection: Connection,
    name: string,
}

export interface QuerySchemasRequest {
    connection: Connection,
    name: string,
}

export interface QueryTablesRequest {
    connection: Connection,
    schema: string,
    name: string,
}

export interface QueryChart {
    name: string,
    value: any,
}

// SPECIFIC (WEB)

export interface QueryRunRequest {
    connection: Connection,
    queryUuid?: string,
    query: string,
    queryOptions?: QueryOptions,
}
