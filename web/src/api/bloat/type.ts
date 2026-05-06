import {Config} from "../database/type"
import {JobStatus} from "./job/type"

// COMMON (WEB AND SERVER)

export interface Bloat {
    uuid: string,
    cluster: string,
    vaultId?: string,
    status: JobStatus,
    command: string,
    commandArgs: string,
    logsPath: string,
    createdAt: string,
}

export interface BloatRequest {
    cluster: string,
    db: Config,
    vaultId?: string,
    target?: BloatTarget,
    options: BloatOptions,
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

// SPECIFIC (WEB)
