import {CSSProperties, ReactNode} from "react";

export interface Response<TData, TError = {}> {
    response: TData,
    error: TError,
}

// NODE
export interface InstanceResponse {
    name: string,
    state: string,
    role: string,
    host: string,
    port: number,
    api_url: string
    lag?: number,
    timeline?: number,
}

export interface Instance extends InstanceResponse{
    api_domain: string,
    leader: boolean,
    inCluster: boolean,
    inInstances: boolean,
}

export interface InstanceMap {
    [instance: string]: Instance,
}

export interface InstanceOverview {
    database_system_identifier: string,
    pending_restart?: boolean
    postmaster_start_time?: string,
    timeline?: number,
    cluster_unlocked?: boolean,
    patroni: { scope: string, version: string },
    state: string,
    role: string,
    xlog?: {
        received_location: number,
        replayed_timestamp: string,
        paused: boolean,
        replayed_location: number
    }
    server_version?: string
}

// CLUSTER
export interface Cluster {
    name: string,
    certId?: string,
    patroniCredId?: string,
    postgresCredId?: string,
    nodes: string[],
    tags?: string[],
}

export interface ClusterMap {
    [name: string]: Cluster,
}

export interface ClusterTabs {
    [key: number]: { body: (cluster: ActiveCluster) => ReactNode, info?: ReactNode }
}

export type DetectionType = "auto" | "manual"

export interface ActiveCluster {
    cluster: Cluster,
    instance: Instance,
    instances: InstanceMap,
    warning: boolean,
    detection: DetectionType,
}

export interface InstanceDetection {
    active: {
        cluster: Cluster,
        instance: Instance,
        instances: InstanceMap,
        warning: boolean,
    },
    colors: ColorsMap,
    fetching: boolean,
    refetch: () => void,
}

// INFO
export interface AppInfo {
    company: string,
    secret: SecretStatus,
}

// SECRET
export interface SecretStatus {
    key: boolean,
    ref: boolean,
}

export interface SecretSetRequest {
    key: string,
    ref: string,
}

export interface SecretUpdateRequest {
    previousKey: string,
    newKey: string,
}

// CREDENTIAL
export interface Credential {
    username: string,
    password: string,
    type: CredentialType,
}

export enum CredentialType {
    POSTGRES,
    PATRONI,
}

export interface CredentialMap {
    [uuid: string]: Credential,
}

// CERT

export interface Cert {
    fileId: string,
    fileName: string,
    path: string,
}

export interface CertMap {
    [uuid: string]: Cert,
}

// BLOAT
export interface Connection {
    host: string,
    port: number,
    credId: string,
}

export interface Target {
    dbName?: string,
    schema?: string,
    table?: string,
    excludeSchema?: string,
    excludeTable?: string,
}

export interface CompactTableRequest {
    cluster: string,
    connection: Connection,
    target?: Target,
    ratio?: number,
}

export interface CompactTable {
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

// COMMON
export interface Style {
    [key: string]: CSSProperties,
}

export interface ColorsMap {
    [name: string]: 'success' | 'primary',
}
