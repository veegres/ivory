import {Alert, AlertColor, AlertTitle, Box, Collapse, InputLabel} from "@mui/material";
import {useState} from "react";
import {AxiosError} from "axios";
import {Style} from "../../app/types";
import {OpenIcon} from "./OpenIcon";

const style: Style = {
    jsonInput: {padding: '10px 0px', whiteSpace: 'pre-wrap'}
}

type ErrorProps = { error: AxiosError | string | unknown, type?: AlertColor }
type GeneralProps = { message: string, type: AlertColor, title?: string, json?: string }

export function ErrorAlert({error, type}: ErrorProps) {
    if (typeof error === "string") return <General type={type ?? "warning"} message={error}/>
    if (!(error instanceof AxiosError)) return <General type={type ?? "warning"} message={"unknown error"}/>
    if (!error.response) return <General type={"error"} message={"ErrorAlert is not detected"}/>

    const {status, statusText} = error.response
    const title = `Error code: ${status ?? 'Unknown'} (${statusText})`
    if (status >= 400 && status < 500) return <General type={"warning"} message={error.message} title={title} json={error.stack}/>
    if (status >= 500) return <General type={"error"} message={error.message} title={title} json={error.stack}/>

    return <General type={"error"} message={error.message} title={title}/>

    function General(props: GeneralProps) {
        const [isOpen, setIsOpen] = useState(false)
        const {message, type, json, title} = props
        const isJson = !!json
        return (
            <Alert severity={type} onClick={() => setIsOpen(!isOpen)} action={<OpenIcon show={isJson} open={isOpen}/>}>
                <AlertTitle>{title ?? type.toString().toUpperCase()}</AlertTitle>
                <Box>Message: {message}</Box>
                <Collapse in={isJson && isOpen}>
                    <InputLabel style={style.jsonInput}>
                        {JSON.stringify(json, null, 4)}
                    </InputLabel>
                </Collapse>
            </Alert>
        )
    }
}
