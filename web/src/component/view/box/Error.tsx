import {Alert, AlertColor, AlertTitle, Box, Collapse, InputLabel} from "@mui/material"
import React, {ReactNode, useState} from "react"

import {SxPropsMap} from "../../../app/type"
import {ClearCacheButton} from "../../widgets/actions/ClearCacheButton"
import {OpenIcon} from "../icon/OpenIcon"

const SX: SxPropsMap = {
    collapse: {display: "flex", flexDirection: "column", gap: 2, marginTop: "20px"},
    label: {fontWeight: "bold"},
    message: {display: "flex", gap: 2, justifyContent: "space-between", alignItems: "center"},
    alert: {"& .MuiAlert-message": {flexGrow: 1}, width: "inherit"},
    input: {whiteSpace: "pre-wrap", wordWrap: "break-word"},
}

type Props = {
    message: string,
    type: AlertColor,
    title?: string,
    stacktrace?: ReactNode,
    params?: string,
    data?: string,
    url?: string,
    clear?: boolean,
}

export function Error(props: Props) {
    const {message, type, stacktrace, title, params, data, url, clear} = props
    const [isOpen, setIsOpen] = useState(false)
    const isStacktrace = !!stacktrace

    return (
        <Alert sx={SX.alert} severity={type} onClick={() => setIsOpen(!isOpen)} action={<OpenIcon show={isStacktrace} open={isOpen}/>}>
            <AlertTitle>{title ?? type.toString().toUpperCase()}</AlertTitle>
            <Box sx={SX.message}>
                <Box>
                    <Box component={"span"}>{message}. </Box>
                    {clear && (
                        <Box component={"span"}>
                            It can be caused because of local cache corruption.
                            You can try to clean it, sometimes it can fix the problem.
                        </Box>
                    )}
                </Box>
                {clear && <ClearCacheButton/>}
            </Box>
            <Collapse in={isStacktrace && isOpen}>
                <Box sx={SX.collapse} onClick={handleDisable}>
                    {renderInfoBox("URL", url)}
                    {renderInfoBox("Request Params", params)}
                    {renderInfoBox("Request Body", data)}
                    {renderInfoBox("Stacktrace", stacktrace)}
                </Box>
            </Collapse>
        </Alert>
    )

    function renderInfoBox(name: string, value: ReactNode) {
        return (
            <Box>
                <Box sx={SX.label}>{name}</Box>
                <InputLabel sx={SX.input}>
                    {JSON.stringify(value, null, 4) ?? "-"}
                </InputLabel>
            </Box>
        )
    }

    function handleDisable(e: React.MouseEvent) {
        e.preventDefault()
        e.stopPropagation()
    }
}
