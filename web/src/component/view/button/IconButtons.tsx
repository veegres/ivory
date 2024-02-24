import {cloneElement, ReactElement} from "react";
import {Box, CircularProgress, IconButton as MuiIconButton, Tooltip} from "@mui/material";
import {
    Add,
    ArrowBack,
    AutoFixHigh,
    Block,
    Cached,
    Cancel,
    CheckCircle,
    ClearAll,
    Close,
    CopyAll,
    Delete,
    Edit,
    Info,
    ManageSearch,
    PendingActions,
    PlayArrow,
    Receipt,
    Restore
} from "@mui/icons-material";

type Color = 'inherit' | 'default' | 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning'
type Placement = "bottom-end" | 'bottom-start' | 'bottom' | 'left-end' | 'left-start' | 'left' | 'right-end' | 'right-start' | 'right' | 'top-end' | 'top-start' | 'top'

type ButtonProps = Props & {
    icon: ReactElement,
    onClick: () => void,
    tooltip: string,
    color?: Color,
    placement?: Placement,
}

export function IconButton(props: ButtonProps) {
    const {icon, color, placement, onClick, tooltip, disabled = false, size = 32} = props
    const loading = props.loading ?? false
    // 0.56 is ratio for size = 32 and fontSize = 18
    const fontSize = Math.floor(size * 0.56)

    return (
        <Tooltip title={tooltip} placement={placement ?? "top"} disableInteractive>
            <Box component={"span"}>
                <MuiIconButton
                    sx={{height: `${size}px`, width: `${size}px`}}
                    color={color}
                    disabled={loading || disabled} onClick={onClick}
                >
                    {loading ? <CircularProgress size={fontSize - 2}/> : cloneElement(icon, {sx: {fontSize}})}
                </MuiIconButton>
            </Box>
        </Tooltip>
    )
}

type Props = {
    loading?: boolean,
    onClick: () => void,
    size?: number,
    disabled?: boolean,
    color?: Color,
    placement?: Placement,
    tooltip?: string,
}

export function EditIconButton(props: Props) {
    const {disabled} = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<Edit/>} tooltip={"Edit"}/>
}

export function DeleteIconButton(props: Props) {
    const {disabled} = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<Delete/>} tooltip={"Delete"}/>
}

export function CancelIconButton(props: Props) {
    const {disabled, tooltip} = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<Cancel/>} tooltip={tooltip ?? "Cancel"}/>
}

export function TerminateIconButton(props: Props) {
    const {disabled} = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<Block/>} tooltip={"Terminate"}/>
}

export function SaveIconButton(props: Props) {
    const {disabled} = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<CheckCircle/>} tooltip={"Save"}/>
}

export function RefreshIconButton(props: Props) {
    const {disabled} = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<Cached/>} tooltip={"Refresh"}/>
}

export function BackIconButton(props: Props) {
    const {disabled} = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<ArrowBack/>} tooltip={"Back"}/>
}

export function CloseIconButton(props: Props) {
    const {disabled} = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<Close/>} tooltip={"Close"}/>
}

export function AddIconButton(props: Props) {
    const {disabled} = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<Add/>} tooltip={"Add"}/>
}

export function AutoIconButton(props: Props) {
    const {disabled} = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<AutoFixHigh/>} tooltip={"Auto"}/>
}

export function PlayIconButton(props: Props) {
    const {disabled} = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<PlayArrow/>} tooltip={"RUN"}/>
}

export function QueryParamsIconButton(props: Props) {
    const {disabled} = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<ManageSearch/>} tooltip={"Query Params"}/>
}

export function RestoreIconButton(props: Props) {
    const {disabled} = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<Restore/>} tooltip={"Restore"}/>
}

export function CopyIconButton(props: Props) {
    const {disabled} = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<CopyAll/>} tooltip={"Copy"}/>
}

export function QueryViewIconButton(props: Props) {
    const {disabled} = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<Receipt/>} tooltip={"Query View"}/>
}

export function LogIconButton(props: Props) {
    const {disabled} = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<PendingActions/>} tooltip={"Log"}/>
}

export function InfoIconButton(props: Props) {
    const {disabled} = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<Info/>} tooltip={"Info"}/>
}

export function ClearAllIconButton(props: Props) {
    const {disabled} = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<ClearAll/>} tooltip={"Clear All"}/>
}


