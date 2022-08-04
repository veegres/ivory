import {cloneElement, ReactElement} from "react";
import {Box, CircularProgress, IconButton, Tooltip} from "@mui/material";
import {Cached, Cancel, CheckCircle, Delete, Edit} from "@mui/icons-material";

const SX = {
    icon: { fontSize: 18 },
    button: { height: "32px", width: "32px" },
}

type ButtonProps = {
    icon: ReactElement,
    loading: boolean,
    onClick: () => void,
    tooltip: string
    disabled?: boolean,
}

function Button({ loading, icon, onClick, tooltip, disabled = false }: ButtonProps) {
    return (
        <Tooltip title={tooltip} placement="top" disableInteractive>
            <Box component={"span"}>
                <IconButton sx={SX.button} disabled={loading || disabled} onClick={onClick}>
                    {loading ? <CircularProgress size={SX.icon.fontSize - 2}/> : cloneElement(icon, {sx: SX.icon})}
                </IconButton>
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
    return <Button icon={<Edit/>} tooltip={"Edit"} loading={loading} disabled={disabled} onClick={onClick}/>
}

export function DeleteIconButton(prop: Props) {
    const { loading, onClick, disabled = false } = prop;
    return <Button icon={<Delete/>} tooltip={"Delete"} loading={loading} disabled={disabled} onClick={onClick}/>
}

export function CancelIconButton(prop: Props) {
    const { loading, onClick, disabled = false } = prop;
    return <Button icon={<Cancel/>} tooltip={"Cancel"} loading={loading} disabled={disabled} onClick={onClick}/>
}

export function SaveIconButton(prop: Props) {
    const { loading, onClick, disabled = false } = prop;
    return <Button icon={<CheckCircle/>} tooltip={"Save"} loading={loading} disabled={disabled} onClick={onClick}/>
}

export function RefreshIconButton(prop: Props) {
    const { loading, onClick, disabled = false } = prop;
    return <Button icon={<Cached/>} tooltip={"Refresh"} loading={loading} disabled={disabled} onClick={onClick}/>
}
