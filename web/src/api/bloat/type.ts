import {Certs} from "../cert/type"
import {Database} from "../postgres"
import {JobStatus} from "./job/type"

// COMMON (WEB AND SERVER)

export interface Connection {
    db: Database,
    certs?: Certs,
    vaultId?: string,
}

export interface BloatTarget {
    database?: string,
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
    connection: Connection,
    target?: BloatTarget,
    options: BloatOptions,
}

export interface Bloat {
    uuid: string,
    status: JobStatus,
    vaultId: string,
    command: string,
    commandArgs: string,
    logsPath: string
}

// SPECIFIC (WEB)
