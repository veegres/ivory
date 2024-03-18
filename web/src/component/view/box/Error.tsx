import {ReactNode, useState} from "react";
import {Alert, AlertColor, AlertTitle, Box, Collapse, InputLabel} from "@mui/material";
import {OpenIcon} from "../icon/OpenIcon";
import {StylePropsMap} from "../../../type/common";

const style: StylePropsMap = {
    stacktrace: {padding: "10px 0px", whiteSpace: "pre-wrap"}
}

type Props = {
    message: string,
    type: AlertColor,
    title?: string,
    stacktrace?: ReactNode,
}

export function Error(props: Props) {
    const {message, type, stacktrace, title} = props
    const [isOpen, setIsOpen] = useState(false)
    const isStacktrace = !!stacktrace

    return (
        <Alert severity={type} onClick={() => setIsOpen(!isOpen)} action={<OpenIcon show={isStacktrace} open={isOpen}/>}>
            <AlertTitle>{title ?? type.toString().toUpperCase()}</AlertTitle>
            <Box>Message: {message}</Box>
            <Collapse in={isStacktrace && isOpen}>
                <InputLabel style={style.stacktrace}>
                    {stacktrace}
                </InputLabel>
            </Collapse>
        </Alert>
    )
}
