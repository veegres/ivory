import {Alert, AlertColor, AlertTitle, Box, Collapse, InputLabel} from "@mui/material";
import {useState} from "react";
import {AxiosError} from "axios";
import {Style} from "../../app/types";
import {OpenIcon} from "./OpenIcon";

const style: Style = {
    jsonInput: {padding: '10px 0px', whiteSpace: 'pre-wrap'}
}

type ErrorProps = { error: AxiosError | string, type?: AlertColor }
type GeneralProps = { message: string, type: AlertColor, title?: string, json?: string }
type JsonProps = { json: string }

export function Error({error, type}: ErrorProps) {
    const [isOpen, setIsOpen] = useState(false)

    if (typeof error === "string") return <General type={type ?? "warning"} message={error}/>
    if (!error.response) return <General type={"error"} message={"Error is not detected"}/>

    const {status, statusText} = error.response
    const title = `Error code: ${status ?? 'Unknown'} (${statusText})`
    if (status >= 400 && status < 500) return <General type={"warning"} message={error.message} title={title} json={error.stack}/>
    if (status >= 500) return <General type={"error"} message={error.message} title={title} json={error.stack}/>

    return <General type={"error"} message={error.message} title={title}/>

    function General(props: GeneralProps) {
        const {message, type} = props
        return (
            <Alert severity={type} onClick={() => setIsOpen(!isOpen)} action={<OpenIcon show={!!props.json} open={isOpen}/>}>
                <AlertTitle>{props.title ?? type.toString().toUpperCase()}</AlertTitle>
                <Box>Message: {message}</Box>
                {props.json ? <Json json={props.json}/> : null}
            </Alert>
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
