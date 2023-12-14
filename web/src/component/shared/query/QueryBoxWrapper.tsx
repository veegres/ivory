import {ReactNode, useState} from "react";
import {Box, SxProps, Theme} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {mergeSxProps} from "../../../app/utils";

const SX: SxPropsMap = {
    box: {padding: "10px", background: "rgba(145,145,145,0.1)", borderRadius: "10px", minHeight: "40px"},
    edit: {outline: 1, outlineColor: "divider"},
    hover: {":hover": {outlineColor: "text.secondary"}},
    focus: {outlineColor: "primary.main"},
}


type Props = {
    editable?: boolean,
    children: ReactNode,
    sx?: SxProps<Theme>,
}

export function QueryBoxWrapper(props: Props) {
    const {children, editable, sx} = props
    const [focus, setFocus] = useState(false)

    const editSx = !editable ? {} : !focus ? mergeSxProps(SX.edit, SX.hover) : mergeSxProps(SX.edit, SX.focus)
    const commonSx = mergeSxProps(SX.box, editSx)
    return (
        <Box sx={mergeSxProps(commonSx, sx)} onFocus={() => setFocus(true)} onBlur={() => setFocus(false)}>
            {children}
        </Box>
    )
}
