import {CSSProperties, ReactElement} from "react";
import {SxProps, Theme} from "@mui/material";
import {SecretStatus} from "./secret";

export interface Response<TData, TError = {}> {
    response: TData,
    error: TError,
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

export interface AppInfo {
    company: string,
    secret: SecretStatus,
}

export interface Connection {
    host: string,
    port: number,
    credId: string,
}

export interface StylePropsMap {
    [key: string]: CSSProperties,
}

export interface SxPropsMap {
    [key: string]: SxProps<Theme>,
}


export interface ColorsMap {
    [name: string]: 'success' | 'primary',
}

export interface EnumOptions {
    name: string,
    label: string,
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
