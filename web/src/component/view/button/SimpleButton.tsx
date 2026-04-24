import {Button} from "@mui/material"
import {ButtonProps} from "@mui/material/Button/Button"

import {SxPropsMap} from "../../../app/type"

const SX: SxPropsMap = {
    button: {padding: "5px", minWidth: 0, borderColor: "divider"},
}

export function SimpleButton(props: ButtonProps) {
    return (
        <Button sx={SX.button} color={"inherit"} variant={"outlined"} {...props}/>
    )
}