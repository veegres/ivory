import {cloneElement, ReactElement} from "react";
import {Box, CircularProgress, IconButton as MuiIconButton, Tooltip} from "@mui/material";
import {ArrowBack, Cached, Cancel, CheckCircle, Close, Delete, Edit} from "@mui/icons-material";

type ButtonProps = Props & {
    icon: ReactElement,
    onClick: () => void,
    tooltip: string
}

export function IconButton(props: ButtonProps) {
    const { icon, onClick, tooltip, disabled = false, size = 32 } = props
    const loading = props.loading ?? false
    // 0.56 is ratio for size = 32 and fontSize = 18
    const fontSize = Math.floor(size * 0.56)

    return (
        <Tooltip title={tooltip} placement={"top"} disableInteractive>
            <Box component={"span"}>
                <MuiIconButton sx={{height: `${size}px`, width: `${size}px`}} disabled={loading || disabled} onClick={onClick}>
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
