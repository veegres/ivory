// COMMON (WEB AND SERVER)

export enum KeeperType {
    PATRONI,
    POSTGRES,
}

export enum Status {
    Active = "ACTIVE",
    Paused = "PAUSED",
}

export type Role = "leader" | "replica" | "unknown";

export interface Keeper {
    host: string,
    port: number,
    name?: string,
    status?: Status,
}

export interface KeeperResponse {
    name?: string,
    status?: Status,
    state: string,
    role: Role,
    lag: number,
    pendingRestart: boolean,
    scheduledSwitchover?: ScheduledSwitchover,
    scheduledRestart?: ScheduledRestart,
    tags?: {[key: string]: any},
    discoveredHost: string,
    discoveredKeeperPort: number,
    discoveredDbPort: number,
}

export interface ScheduledSwitchover {
    at: string,
    to: string,
}

export interface ScheduledRestart {
    at: string,
    pendingRestart: boolean,
}
