import {Connection} from "./common";

export interface BloatTarget {
    dbName?: string,
    schema?: string,
    table?: string,
    excludeSchema?: string,
    excludeTable?: string,
}

export interface BloatRequest {
    cluster: string,
    connection: Connection,
    target?: BloatTarget,
    ratio?: number,
}

export interface Bloat {
    uuid: string,
    status: JobStatus,
    credentialId: string,
    command: string,
    commandArgs: string,
    logsPath: string
}

export enum JobStatus {
    PENDING,
    RUNNING,
    FINISHED,
    FAILED,
    STOPPED,
    UNKNOWN,
}

export enum EventStream {
    START = "start",
    END = "end",
}
