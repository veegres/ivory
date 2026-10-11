import {useSnackbar} from "../../provider/SnackbarProvider"
import {CopyIconButton} from "./IconButtons"

type Props = {
    value: string,
    tooltip?: string,
    message?: string,
}

export function CopyButton(props: Props) {
    const {value, tooltip, message} = props
    const snackbar = useSnackbar()

    return <CopyIconButton tooltip={tooltip ?? "Copy to your clipboard"} onClick={handleCopy}/>

    function handleCopy() {
        navigator.clipboard.writeText(value)
            .then(() => snackbar(message ?? "Copied to clipboard!", "info"))
            .catch(() => snackbar("Failed to copy to clipboard", "error"))
    }
}
