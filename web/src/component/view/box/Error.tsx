import {ReactNode, useState} from "react";
import {Alert, AlertColor, AlertTitle, Box, Collapse, InputLabel} from "@mui/material";
import {OpenIcon} from "../icon/OpenIcon";
import {StylePropsMap, SxPropsMap} from "../../../type/general";

const SX: SxPropsMap = {
    collapse: {display: "flex", flexDirection: "column", gap: 2, marginTop: "20px"},
    label: {fontWeight: "bold"},
}

const style: StylePropsMap = {
    input: {whiteSpace: "pre-wrap"},
}

type Props = {
    message: string,
    type: AlertColor,
    title?: string,
    stacktrace?: ReactNode,
    params?: string,
    data?: string,
    url?: string,
}

export function Error(props: Props) {
    const {message, type, stacktrace, title, params, data, url} = props
    const [isOpen, setIsOpen] = useState(false)
    const isStacktrace = !!stacktrace

    return (
        <Alert severity={type} onClick={() => setIsOpen(!isOpen)} action={<OpenIcon show={isStacktrace} open={isOpen}/>}>
            <AlertTitle>{title ?? type.toString().toUpperCase()}</AlertTitle>
            <Box>Message: {message}</Box>
            <Collapse in={isStacktrace && isOpen}>
                <Box sx={SX.collapse}>
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
                <InputLabel style={style.input}>
                    {JSON.stringify(value, null, 4) ?? "-"}
                </InputLabel>
            </Box>
        )
    }
}
