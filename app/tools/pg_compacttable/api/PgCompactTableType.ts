import {DbConfig} from "../../../features/query/api/QueryType"
import {JobStatus} from "./job/PgCompactTableJobType"

// COMMON (WEB AND SERVER)

export interface PgCompactTable {
    uuid: string,
    cluster: string,
    vaultId?: string,
    status: JobStatus,
    command: string,
    commandArgs: string,
    logsPath: string,
    createdAt: string,
}

export interface PgCompactTableRequest {
    cluster: string,
    db: DbConfig,
    vaultId?: string,
    target?: PgCompactTableTarget,
    options: PgCompactTableOptions,
}

export interface PgCompactTableTarget {
    database?: string,
    schema?: string,
    table?: string,
    excludeSchema?: string,
    excludeTable?: string,
}

export interface PgCompactTableOptions {
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
