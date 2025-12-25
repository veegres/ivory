import {ReactNode} from "react"

import {Certs} from "../cert/type"
import {
    InstanceScheduledRestart,
    InstanceScheduledSwitchover, Role,
    Sidecar,
} from "../instance/type"
import {Permission} from "../permission/type"
import {Database} from "../postgres"

// COMMON (WEB AND SERVER)

export interface ClusterOptions {
    tls: ClusterTls,
    certs: Certs,
    credentials: Credentials,
    tags?: string[],
}

export interface Cluster extends ClusterOptions {
    name: string,
    sidecars: Sidecar[],
    sidecarsOverview: InstanceOverview,
}

export interface ClusterAuto extends ClusterOptions {
    name: string,
    instance: Sidecar,
}

export interface ClusterTls {
    sidecar: boolean,
    database: boolean,
}

export interface Credentials {
    patroniId?: string,
    postgresId?: string,
}


export interface ClusterOverview {
    instances: InstanceOverview,
    detectedBy?: Sidecar,
    mainInstance?: Instance,
}

export interface InstanceOverview {
    [sidecar: string]: Instance | undefined
}

export interface Instance {
    state: string,
    role: Role,
    lag: number,
    pendingRestart: boolean,
    database: Database,
    sidecar: Sidecar,
    scheduledSwitchover?: InstanceScheduledSwitchover,
    scheduledRestart?: InstanceScheduledRestart,
    tags?: {[key: string]: any},
    inCluster: boolean,
    inSidecar: boolean,
}

// SPECIFIC (WEB)

export interface ActiveCluster {
    cluster: Cluster,
    detectBy?: Instance,
    warning: boolean,
}

export interface ActiveInstance {
    [cluster: string]: string | undefined
}

export interface ClusterTab {
    label: string,
    body: (cluster: Cluster, overview?: ClusterOverview) => ReactNode,
    permission: Permission,
    info?: ReactNode,
}
