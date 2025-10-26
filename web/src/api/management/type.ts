import {AuthType} from "../auth/type"
import {AvailabilityConfig} from "../config/type"
import {SecretStatus} from "../secret/type"

// COMMON (WEB AND SERVER)

export interface R<TData, TError = object> {
    response: TData,
    error: TError,
}


export interface AppInfo {
    auth: AuthInfo,
    config: ConfigInfo,
    secret: SecretStatus,
    version: Version,
}

export interface ConfigInfo {
    company: string,
    configured: boolean,
    availability: AvailabilityConfig,
    error?: string,
}


export interface AuthInfo {
    supported: AuthType[],
    authorised: boolean,
    error?: string,
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

