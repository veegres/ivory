import {FileUsageType} from "./common";

export interface Cert {
    fileName: string,
    fileUsageType: FileUsageType,
    path: string,
    type: CertType,
}

export interface CertMap {
    [uuid: string]: Cert,
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
