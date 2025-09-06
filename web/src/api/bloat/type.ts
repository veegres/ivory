import {DbConnection} from "../management/type";
import {JobStatus} from "./job/type";

// COMMON (WEB AND SERVER)

export interface BloatTarget {
    dbName?: string,
    schema?: string,
    table?: string,
    excludeSchema?: string,
    excludeTable?: string,
}

export interface BloatOptions {
    force?: boolean,
    noReindex?: boolean,
    noInitialVacuum?: boolean,
    initialReindex?: boolean,
    routineVacuum?: boolean,
    delayRatio?: number,
    minTableSize?: number,
    maxTableSize?: number,
}

export interface BloatRequest {
    cluster: string,
    connection: DbConnection,
    target?: BloatTarget,
    options: BloatOptions,
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
