import {AuthType} from "../auth/type"
import {PermissionStatus} from "../permission/type"
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
    error?: string,
}

export interface AuthInfo {
    supported: AuthType[],
    authorised: boolean,
    user?: UserInfo,
    error?: string,
}

export interface UserInfo {
    username: string,
    permissions?: { [key: string]: PermissionStatus },
}

export interface SecretUpdateRequest {
    previousKey: string,
    newKey: string,
}

// SPECIFIC (WEB)

export interface Version {
    tag: string,
    commit: string,
    label: string,
}

