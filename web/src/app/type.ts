import {Theme} from "@mui/material"
import {SvgIconProps} from "@mui/material/SvgIcon/SvgIcon"
import {SystemStyleObject} from "@mui/system/styleFunctionSx/styleFunctionSx"
import {CSSProperties, ReactElement} from "react"

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