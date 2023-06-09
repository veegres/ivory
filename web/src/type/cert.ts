import {FileUsageType} from "./common";
import {AxiosProgressEvent} from "axios";

// COMMON (WEB AND SERVER)

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

export interface CertAddRequest {
    path: string,
    type: CertType,
}

// SPECIFIC (WEB)

export interface CertUploadRequest {
    file: File,
    type: CertType,
    setProgress?: (progressEvent: AxiosProgressEvent) => void
}

export interface CertTabs {
    [key: number]: { label: string, type: CertType }
}
