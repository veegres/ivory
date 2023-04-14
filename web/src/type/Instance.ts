import {Database, Sidecar} from "./common";
import {ReactNode} from "react";
import {Certs} from "./cluster";

export enum InstanceTabType {QUERY, CHART}
export interface InstanceTab {
    label: string,
    body: (cluster: string, db: Database) => ReactNode,
    info?: ReactNode,
}

export interface InstanceRequest {
    credentialId?: string,
    certs: Certs,
    host: string,
    port: number,
    body?: any,
}

export interface InstanceInfo {
    state: string,
    role: string,
    sidecar: Sidecar,
}


export interface InstanceResponse {
    state: string,
    role: string,
    lag: number,
    database: Database,
    sidecar: Sidecar,
}

export interface Instance extends InstanceResponse {
    leader: boolean,
    inCluster: boolean,
    inInstances: boolean,
}

export interface InstanceMap {
    [instance: string]: Instance,
}
