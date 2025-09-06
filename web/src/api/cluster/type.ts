import {ReactNode} from "react";
import {ColorsMap} from "../management/type";
import {InstanceMap, InstanceWeb, Sidecar} from "../instance/type";
import {Certs} from "../cert/type";

// COMMON (WEB AND SERVER)

export interface ClusterOptions {
    tls: ClusterTls,
    certs: Certs,
    credentials: Credentials,
    tags?: string[],
}

export interface Cluster extends ClusterOptions {
    name: string,
    instances: Sidecar[],
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

// SPECIFIC (WEB)

// TODO should we return map from backend?
export interface ClusterMap {
    [name: string]: Cluster,
}

export interface ActiveCluster {
    cluster: Cluster,
    defaultInstance: InstanceWeb,
    combinedInstanceMap: InstanceMap,
    warning: boolean,
    detection: DetectionType,
}

export interface ClusterTabs {
    [key: number]: {
        label: string,
        body: (cluster: ActiveCluster, activeInstance?: InstanceWeb) => ReactNode,
        info?: ReactNode,
    }
}

export type DetectionType = "auto" | "manual"

export interface InstanceDetection {
    defaultInstance: InstanceWeb,
    combinedInstanceMap: InstanceMap,
    detection: DetectionType,
    warning: boolean,
    colors: ColorsMap,
    fetching: boolean,
    active: boolean,
    refetch: () => void,
}

