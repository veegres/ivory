import {CSSProperties, ReactElement} from "react";
import {SxProps, Theme} from "@mui/material";
import {SecretStatus} from "./secret";

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
    credId: string,
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
}

export enum FileUsageType {
    UPLOAD,
    PATH,
}

export interface Database {
    host: string,
    port: number,
    database?: string,
}

export interface Sidecar {
    host: string,
    port: number,
}

// SPECIFIC (WEB)

export interface StylePropsMap {
    [key: string]: CSSProperties,
}

export interface SxPropsMap {
    [key: string]: SxProps<Theme>,
}


export interface ColorsMap {
    [name: string]: "success" | "primary" | "error" | "warning",
}

export interface EnumOptions {
    label: string,
    key: string,
    name?: string,
    icon?: ReactElement,
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
