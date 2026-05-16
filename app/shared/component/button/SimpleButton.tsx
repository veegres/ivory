import {Button} from "@mui/material"
import {ButtonProps} from "@mui/material"

import {SxPropsMap} from "../../helper/type"

const SX: SxPropsMap = {
    button: {padding: "5px", minWidth: 0, borderColor: "divider"},
}

export function SimpleButton(props: ButtonProps) {
    return (
        <Button color={"inherit"} variant={"outlined"} {...props} sx={{...SX.button, ...props.sx}}/>
    )
}