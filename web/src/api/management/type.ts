import {CSSProperties, ReactElement} from "react";
import {Theme} from "@mui/material";
import {SecretStatus} from "../secret/type";
import {SystemStyleObject} from "@mui/system/styleFunctionSx/styleFunctionSx";
import {SvgIconProps} from "@mui/material/SvgIcon/SvgIcon";
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
    icon: ReactElement<SvgIconProps>,
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
