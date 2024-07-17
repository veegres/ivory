import {CSSProperties, ReactElement} from "react";
import {Theme} from "@mui/material";
import {SecretStatus} from "./secret";
import {SystemStyleObject} from "@mui/system/styleFunctionSx/styleFunctionSx";

// COMMON (WEB AND SERVER)

export interface Response<TData, TError = object> {
    response: TData,
    error: TError,
}

export interface Login {
    username: string,
    password: string,
}

export interface DbConnection {
    host: string,
    port: number,
    credentialId: string,
}

export enum AuthType {
    NONE,
    BASIC,
}

export interface AppInfo {
    company: string,
    configured: boolean,
    secret: SecretStatus,
    availability: Availability,
    version: Version,
    auth: AuthInfo,
}

export interface AuthInfo {
    authorised: boolean,
    error: string,
    type: AuthType,
}

export interface AppConfig {
    company: string,
    availability: Availability,
    auth: AuthConfig,
}

export interface AuthConfig {
    type: AuthType,
    body: any,
}

export interface Availability {
    manualQuery: boolean,
}

export interface Version {
    tag: string,
    commit: string,
    label: string,
}

export enum FileUsageType {
    UPLOAD,
    PATH,
}

export interface Database {
    host: string,
    port: number,
    name?: string,
}

export enum SidecarStatus {
    Active = "ACTIVE",
    Paused = "PAUSED",
}
export interface Sidecar {
    host: string,
    port: number,
    status?: SidecarStatus,
}

// SPECIFIC (WEB)

export interface StylePropsMap {
    [key: string]: CSSProperties,
}

export interface SxPropsMap {
    // NOTE: we don't need usage of array here that is why it is not use SxProps<Theme> interface
    [key: string]: SystemStyleObject<Theme> | ((theme: Theme) => SystemStyleObject<Theme>),
}

export interface ColorsMap {
    [name: string]: "success" | "primary" | "error" | "warning",
}

export interface EnumOptions {
    label: string,
    key: string,
    name?: string,
    icon: ReactElement,
    color?: string,
    badge?: string,
}

export enum Settings {
    MENU,
    PASSWORD,
    CERTIFICATE,
    SECRET,
    ABOUT,
}

export interface Links {
    [key: string]: {
        name: string,
        link: string,
    }
}
