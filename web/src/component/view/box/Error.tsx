import {useState} from "react";
import {Alert, AlertColor, AlertTitle, Box, Collapse, InputLabel} from "@mui/material";
import {OpenIcon} from "../icon/OpenIcon";
import {StylePropsMap} from "../../../type/common";

const style: StylePropsMap = {
    jsonInput: {padding: '10px 0px', whiteSpace: 'pre-wrap'}
}

type Props = {
    message: string,
    type: AlertColor,
    title?: string,
    json?: string
}

export function Error(props: Props) {
    const {message, type, json, title} = props
    const [isOpen, setIsOpen] = useState(false)
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
