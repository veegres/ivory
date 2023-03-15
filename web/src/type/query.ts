import {Database} from "./common";

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

export interface QueryRunRequest {
    queryUuid: string,
    clusterName: string,
    db: Database,
}

export interface QueryKillRequest {
    pid: number,
    clusterName: string,
    db: Database,
}

export interface QueryChartRequest {
    clusterName: string,
    db: Database,
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
