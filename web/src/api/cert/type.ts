import {AxiosProgressEvent} from "axios";

// COMMON (WEB AND SERVER)

export enum FileUsageType {
    UPLOAD,
    PATH,
}

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

export interface Certs {
    clientCAId?: string,
    clientCertId?: string,
    clientKeyId?: string,
}