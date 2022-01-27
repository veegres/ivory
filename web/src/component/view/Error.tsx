import {Alert, AlertColor, AlertTitle, Box, Collapse, IconButton, InputLabel} from "@mui/material";
import {useState} from "react";
import {KeyboardArrowDown, KeyboardArrowUp} from "@mui/icons-material";
import {AxiosError} from "axios";
import {Style} from "../../app/types";


const style: Style = {
    jsonInput: { padding: '10px 0px', whiteSpace: 'pre-wrap' }
}

type ErrorProps = { error: AxiosError | string }
type GeneralProps = { message: string, type: AlertColor, title?: string, json?: string }
type JsonProps = { json: string }
type Button = { isShown: boolean }

export function Error({ error }: ErrorProps) {
    const [isOpen, setIsOpen] = useState(false)

    if (typeof error === "string") return <General type={"error"} message={error} />
    if (!error.response) return <General type={"error"} message={"Error is not detected"} />

    const { status, statusText } = error.response
    const title = `Error code: ${status ?? 'Unknown'} (${statusText})`
    if (status >= 400 && status < 500) return <General type={"warning"} message={error.message} title={title} json={error.stack} />
    if (status >= 500) return <General type={"error"} message={error.message} title={title} json={error.stack} />

    return <General type={"error"} message={error.message} title={title} />

    function General(props: GeneralProps) {
        const { message, type } = props
        return (
            <Alert severity={type} onClick={() => setIsOpen(!isOpen)} action={<Button isShown={!!props.json} />}>
                <AlertTitle>{props.title ?? type.toString()}</AlertTitle>
                <Box>Message: {message}</Box>
                {props.json ? <Json json={props.json} /> : null}
            </Alert>
        )
    }

    function Button(props: Button) {
        if (!props.isShown) return null

        return (
            <IconButton color={"inherit"} disableRipple>
                {isOpen ? <KeyboardArrowUp /> : <KeyboardArrowDown />}
            </IconButton>
        )
    }

    function Json(props: JsonProps) {
        return (
            <Collapse in={isOpen}>
                <InputLabel style={style.jsonInput}>
                    {JSON.stringify(props.json, null, 4)}
                </InputLabel>
            </Collapse>
        )
    }
}
