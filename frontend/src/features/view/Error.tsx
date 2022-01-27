import {Alert, AlertColor, AlertTitle, Box, Collapse, IconButton, InputLabel} from "@mui/material";
import {useState} from "react";
import {KeyboardArrowDown, KeyboardArrowUp} from "@mui/icons-material";
import {AxiosError} from "axios";

export function Error({ error }: { error: AxiosError }) {
    const [isOpen, setIsOpen] = useState(false)

    if (!error.response) return <General type={"error"} message={"Error is not detected"} />

    const { status, statusText } = error.response
    switch (status) {
        case 404: return <General type={"warning"} message={error.message} />
        default: return <General type={"error"} message={error.message} />
    }

    function General(props: { message: string, type: AlertColor }) {
        const { message, type } = props
        return (
            <Alert severity={type} onClick={() => setIsOpen(!isOpen)} action={<Button />}>
                <AlertTitle>Error code: {status ?? 'Unknown'} ({statusText})</AlertTitle>
                <Box>Message: {message}</Box>
                <Stacktrace />
            </Alert>
        )
    }

    function Button() {
        return (
            <IconButton color={"inherit"} disableRipple>
                {isOpen ? <KeyboardArrowUp /> : <KeyboardArrowDown />}
            </IconButton>
        )
    }

    function Stacktrace() {
        return (
            <Collapse in={isOpen}>
                <InputLabel sx={{ padding: '10px 0px', whiteSpace: 'pre-wrap' }}>
                    {!error?.stack ? 'None' : JSON.stringify(error.stack, null, 4)}
                </InputLabel>
            </Collapse>
        )
    }
}
