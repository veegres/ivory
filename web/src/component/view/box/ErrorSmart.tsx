import {AlertColor} from "@mui/material";
import {AxiosError} from "axios";
import {Error} from "./Error";

type Props = {
    error: AxiosError | string | unknown,
    type?: AlertColor
}

export function ErrorSmart({error, type}: Props) {
    if (typeof error === "string") return <Error type={type ?? "warning"} message={error}/>
    if (!(error instanceof AxiosError)) return <Error type={type ?? "warning"} message={String(error)}/>
    if (!error.response) return <Error type={"error"} message={"ErrorAlert is not detected"}/>

    const {status, statusText} = error.response
    const title = `Error code: ${status ?? 'Unknown Code'} (${statusText ?? 'Unknown Error Name'})`
    const json = JSON.stringify(error.stack, null, 4)
    if (status >= 400 && status < 500) return <Error type={"warning"} message={error.response.data.error} title={title} stacktrace={json}/>
    if (status >= 500) return <Error type={"error"} message={error.message} title={title} stacktrace={json}/>

    return <Error type={"error"} message={error.message} title={title}/>
}
