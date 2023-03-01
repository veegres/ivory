import {CSSProperties, ReactElement, ReactNode} from "react";
import {SxProps, Theme} from "@mui/material";

export interface Response<TData, TError = {}> {
    response: TData,
    error: TError,
}

// INSTANCE
export interface InstanceRequest {
    cluster: string,
    host: string,
    port: number,
    body?: any,
}

export interface Instance {
    state: string,
    role: string,
    lag: number,
    database: Database,
    sidecar: Sidecar,
}

export interface Database {
    host: string,
    port: number,
}

export interface Sidecar {
    host: string,
    port: number,
}

export interface DefaultInstance extends Instance {
    leader: boolean,
    inCluster: boolean,
    inInstances: boolean,
}

export interface InstanceMap {
    [instance: string]: DefaultInstance,
}

export interface InstanceInfo {
    state: string,
    role: string,
    sidecar: Sidecar,
}

// CLUSTER
export interface Cluster {
    name: string,
    certs: Certs,
    credentials: Credentials,
    instances: Sidecar[],
    tags?: string[],
}

export interface Certs {
    clientCAId?: string,
    clientCertId?: string,
    clientKeyId?: string,
}

export interface Credentials {
    patroniId?: string,
    postgresId?: string,
}

export interface ClusterMap {
    [name: string]: Cluster,
}

export interface ClusterTabs {
    [key: number]: {
        label: string,
        body: (cluster: ActiveCluster) => ReactNode,
        info?: ReactNode,
    }
}

export type DetectionType = "auto" | "manual"

export interface ActiveCluster {
    cluster: Cluster,
    defaultInstance: DefaultInstance,
    combinedInstanceMap: InstanceMap,
    warning: boolean,
    detection: DetectionType,
}

export interface ActiveInstance {
    cluster: string,
    sidecar: Sidecar,
    database: Database,
}

export interface InstanceDetection {
    defaultInstance: DefaultInstance,
    combinedInstanceMap: InstanceMap,
    detection: DetectionType,
    warning: boolean,
    colors: ColorsMap,
    fetching: boolean,
    active: boolean,
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
    fileName: string,
    fileUsageType: FileUsageType,
    path: string,
    type: CertType,
}

export interface CertMap {
    [uuid: string]: Cert,
}

export enum FileUsageType {
    UPLOAD,
    PATH,
}

export enum CertType {
    CLIENT_CA,
    CLIENT_CERT,
    CLIENT_KEY,
}

export interface CertUploadRequest {
    file: File,
    type: CertType,
    setProgress?: (progressEvent: ProgressEvent) => void
}

export interface CertAddRequest {
    path: string,
    type: CertType,
}

export interface CertTabs {
    [key: number]: { label: string, type: CertType }
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

// QUERY
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

export interface QueryMap {
    [key: string]: Query,
}

export interface QueryRunRequest {
    queryUuid: string,
    clusterName: string,
    db: Database,
}


export interface QueryField {
    name: string,
    dataType: string,
    dataTypeOID: string,
}

export interface QueryRunResponse {
    fields: QueryField[],
    rows: any[][],
}

// SETTINGS
export enum Settings {
    MENU,
    PASSWORD,
    CERTIFICATE,
    SECRET,
    ABOUT,
}

// COMMON
export interface StylePropsMap {
    [key: string]: CSSProperties,
}

export interface SxPropsMap {
    [key: string]: SxProps<Theme>,
}

export interface ColorsMap {
    [name: string]: 'success' | 'primary',
}

export interface EnumOptions {
    name: string,
    label: string,
    icon: ReactElement,
    color?: string,
    badge?: string,
}
