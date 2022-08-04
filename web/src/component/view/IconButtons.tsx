import {cloneElement, ReactElement} from "react";
import {Box, CircularProgress, IconButton as MuiIconButton, Tooltip} from "@mui/material";
import {Cached, Cancel, CheckCircle, Delete, Edit} from "@mui/icons-material";

type ButtonProps = {
    icon: ReactElement,
    size?: number,
    loading: boolean,
    onClick: () => void,
    tooltip: string
    disabled?: boolean,
}

export function IconButton({ loading, icon, onClick, tooltip, disabled = false, size = 32 }: ButtonProps) {
    const fontSize = Math.floor(size * 0.56) // 0.56 is ratio for size = 32 and fontSize = 18
    return (
        <Tooltip title={tooltip} placement="top" disableInteractive>
            <Box component={"span"}>
                <MuiIconButton sx={{ height: `${size}px`, width: `${size}px` }} disabled={loading || disabled} onClick={onClick}>
                    {loading ? <CircularProgress size={fontSize - 2}/> : cloneElement(icon, {sx: { fontSize }})}
                </MuiIconButton>
            </Box>
        </Tooltip>
    )
}

type Props = {
    loading: boolean,
    onClick: () => void,
    disabled?: boolean,
}

export function EditIconButton(prop: Props) {
    const { loading, onClick, disabled = false } = prop;
    return <IconButton icon={<Edit/>} tooltip={"Edit"} loading={loading} disabled={disabled} onClick={onClick}/>
}

export function DeleteIconButton(prop: Props) {
    const { loading, onClick, disabled = false } = prop;
    return <IconButton icon={<Delete/>} tooltip={"Delete"} loading={loading} disabled={disabled} onClick={onClick}/>
}

export function CancelIconButton(prop: Props) {
    const { loading, onClick, disabled = false } = prop;
    return <IconButton icon={<Cancel/>} tooltip={"Cancel"} loading={loading} disabled={disabled} onClick={onClick}/>
}

export function SaveIconButton(prop: Props) {
    const { loading, onClick, disabled = false } = prop;
    return <IconButton icon={<CheckCircle/>} tooltip={"Save"} loading={loading} disabled={disabled} onClick={onClick}/>
}

export function RefreshIconButton(prop: Props) {
    const { loading, onClick, disabled = false } = prop;
    return <IconButton icon={<Cached/>} tooltip={"Refresh"} loading={loading} disabled={disabled} onClick={onClick}/>
}
