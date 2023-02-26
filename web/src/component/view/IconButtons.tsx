import {cloneElement, ReactElement} from "react";
import {Box, CircularProgress, IconButton as MuiIconButton, Tooltip} from "@mui/material";
import {
    Add, ArrowBack, AutoFixHigh, Cached, Cancel, CheckCircle, Close, CopyAll, Delete, Edit, PlayArrow, Undo
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
        <Tooltip title={tooltip} placement={placement} disableInteractive>
            <Box component={"span"}>
                <MuiIconButton sx={{height: `${size}px`, width: `${size}px`}} color={color} disabled={loading || disabled} onClick={onClick}>
                    {loading ? <CircularProgress size={fontSize - 2}/> : cloneElement(icon, {sx: { fontSize }})}
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
}

export function EditIconButton(props: Props) {
    const { disabled } = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<Edit/>} tooltip={"Edit"}/>
}

export function DeleteIconButton(props: Props) {
    const { disabled } = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<Delete/>} tooltip={"Delete"}/>
}

export function CancelIconButton(props: Props) {
    const { disabled } = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<Cancel/>} tooltip={"Cancel"}/>
}

export function SaveIconButton(props: Props) {
    const { disabled } = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<CheckCircle/>} tooltip={"Save"} />
}

export function RefreshIconButton(props: Props) {
    const { disabled } = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<Cached/>} tooltip={"Refresh"}/>
}

export function BackIconButton(props: Props) {
    const { disabled } = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<ArrowBack/>} tooltip={"Back"}/>
}

export function CloseIconButton(props: Props) {
    const { disabled } = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<Close/>} tooltip={"Close"}/>
}

export function AddIconButton(props: Props) {
    const { disabled } = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<Add/>} tooltip={"Add"}/>
}

export function AutoIconButton(props: Props) {
    const { disabled } = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<AutoFixHigh/>} tooltip={"Auto Detection (coming soon)"}/>
}

export function PlayIconButton(props: Props) {
    const { disabled } = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<PlayArrow/>} tooltip={"Run"}/>
}

export function UndoIconButton(props: Props) {
    const { disabled } = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<Undo/>} tooltip={"Undo"}/>
}

export function CopyIconButton(props: Props) {
    const { disabled } = props
    return <IconButton {...props} disabled={disabled ?? false} icon={<CopyAll/>} tooltip={"Copy"}/>
}

