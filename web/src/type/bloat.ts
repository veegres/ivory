import {DbConnection} from "./general";
import {JobStatus} from "./job";

// COMMON (WEB AND SERVER)

export interface BloatTarget {
    dbName?: string,
    schema?: string,
    table?: string,
    excludeSchema?: string,
    excludeTable?: string,
}

export interface BloatRequest {
    cluster: string,
    connection: DbConnection,
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

// SPECIFIC (WEB)
