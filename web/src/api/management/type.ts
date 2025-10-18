import {SecretStatus} from "../secret/type";
import {AuthInfo} from "../auth/type";
import {Availability} from "../config/type";

// COMMON (WEB AND SERVER)

export interface R<TData, TError = object> {
    response: TData,
    error: TError,
}

export interface AppInfo {
    company: string,
    configured: boolean,
    secret: SecretStatus,
    availability: Availability,
    version: Version,
    auth: AuthInfo,
}

export interface SecretUpdateRequest {
    previousKey: string,
    newKey: string,
}

// ENV VARIABLES (WEB AND SERVER)

export interface Version {
    tag: string,
    commit: string,
    label: string,
}

// SPECIFIC (WEB)

