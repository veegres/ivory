import {Database, Sidecar} from "./common";
import {ReactNode} from "react";
import {Certs} from "./cluster";

// COMMON (WEB AND SERVER)

export interface InstanceRequest {
    credentialId?: string,
    certs: Certs,
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

export interface InstanceInfo {
    state: string,
    role: string,
    sidecar: Sidecar,
}

// SPECIFIC (WEB)

// TODO should we return map from the server?
export interface InstanceMap {
    [instance: string]: InstanceWeb,
}

export enum InstanceTabType {QUERY, CHART}
export interface InstanceTab {
    label: string,
    body: (cluster: string, db: Database) => ReactNode,
    info?: ReactNode,
}

export interface InstanceWeb extends Instance {
    leader: boolean,
    inCluster: boolean,
    inInstances: boolean,
}
