import {ReactNode, useState} from "react";
import {Box} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {mergeSxProps} from "../../../app/utils";
import {blue} from "@mui/material/colors";

const SX: SxPropsMap = {
    box: {padding: "10px", background: "rgba(145,145,145,0.1)", borderRadius: "10px"},
    edit: {border: 1, borderColor: "divider"},
    hover: {":hover": {borderColor: "text.secondary"}},
    focus: {borderColor: "primary.main", boxShadow: `inset 0px 0px 0px 1px ${blue[500]}`},
}


type Props = {
    editable?: boolean,
    children: ReactNode,
}

export function QueryDescription(props: Props) {
    const {children, editable} = props
    const [focus, setFocus] = useState(false)

    const editProps = !editable ? {} : !focus ? mergeSxProps(SX.edit, SX.hover) : mergeSxProps(SX.edit, SX.focus)
    return (
        <Box sx={mergeSxProps(SX.box, editProps)} onFocus={() => setFocus(true)} onBlur={() => setFocus(false)}>
            {children}
        </Box>
    )
}
