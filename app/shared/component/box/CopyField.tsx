import {Box} from "@mui/material"

import {SxPropsMap} from "../../helper/HelperType"
import {useSnackbar} from "../../provider/SnackbarProvider"
import {CopyIconButton} from "../button/IconButtons"

const SX: SxPropsMap = {
    field: {
        display: "flex", alignItems: "center", gap: 1, height: "40px",
        padding: "0px 4px 0px 12px", border: 1, borderColor: "divider", borderRadius: 1,
    },
    value: {
        flexGrow: 1, minWidth: 0, fontFamily: "monospace", fontSize: "12px", color: "text.secondary",
        whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis",
    },
}

type Props = {
    value: string,
    tooltip?: string,
    message?: string,
}

// CopyField is a single-line value that exists to be copied somewhere else - an
// endpoint, a token, a registration link. It is a field rather than plain text
// because the copy button is the point: nobody should be retyping these by
// hand. Several lines of code belong in a CodeField instead.
export function CopyField(props: Props) {
    const {value, tooltip, message} = props
    const snackbar = useSnackbar()

    return (
        <Box sx={SX.field}>
            <Box sx={SX.value} title={value}>{value}</Box>
            <CopyIconButton tooltip={tooltip ?? "Copy to your clipboard"} onClick={handleCopy}/>
        </Box>
    )

    function handleCopy() {
        navigator.clipboard.writeText(value).then(() => {
            snackbar(message ?? "Copied to clipboard!", "info")
        })
    }
}
