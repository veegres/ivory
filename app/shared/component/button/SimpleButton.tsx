import {Box, Button, ButtonProps, Tooltip} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../helper/HelperType"

type Placement = "bottom-end" | "bottom-start" | "bottom" | "left-end" | "left-start" | "left" | "right-end" | "right-start" | "right" | "top-end" | "top-start" | "top"

const SX: SxPropsMap = {
    button: {padding: "3px", minWidth: 0, borderColor: "divider"},
}

type Props = ButtonProps & {
    tooltip?: ReactNode,
    placement?: Placement,
    arrow?: boolean,
}

export function SimpleButton(props: Props) {
    const {tooltip, placement, arrow = true, ...buttonProps} = props
    const button = <Button color={"inherit"} variant={"outlined"} {...buttonProps} sx={{...SX.button, ...buttonProps.sx}}/>

    if (!tooltip) return button

    return (
        <Tooltip title={tooltip} placement={placement ?? "top"} arrow={arrow} disableInteractive>
            <Box component={"span"}>{button}</Box>
        </Tooltip>
    )
}
