import {AlertColor} from "@mui/material";
import {AxiosError} from "axios";
import {Error} from "./Error";

type Props = {
    error: AxiosError | string | unknown,
    type?: AlertColor
}

export function ErrorSmart({error, type}: Props) {
    if (typeof error === "string") return <Error type={type ?? "warning"} message={error}/>
    if (!(error instanceof AxiosError)) return <Error type={type ?? "warning"} message={"unknown error"}/>
    if (!error.response) return <Error type={"error"} message={"ErrorAlert is not detected"}/>

    const {status, statusText} = error.response
    const title = `Error code: ${status ?? 'Unknown Code'} (${statusText ?? 'Unknown Error Name'})`
    if (status >= 400 && status < 500) return <Error type={"warning"} message={error.response.data.error} title={title} json={error.stack}/>
    if (status >= 500) return <Error type={"error"} message={error.message} title={title} json={error.stack}/>

    return <Error type={"error"} message={error.message} title={title}/>
}
