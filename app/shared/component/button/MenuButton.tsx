import {Box, Collapse} from "@mui/material"
import {ReactNode, useState} from "react"

import {SxPropsMap} from "../../helper/HelperType"
import {MoreIconButton} from "./IconButtons"

const SX: SxPropsMap = {
    box: {position: "relative"},
    collapse: {position: "absolute", top: "50%", transform: "translate(calc(-100% + -2px), -50%)"},
    paper: {
        display: "flex", gap: "3px", alignItems: "center", padding: "2px",
        border: 1, borderRadius: 1, borderColor: "divider", backgroundColor: "background.default",
    },
}

type Props = {
    children: ReactNode,
    open?: boolean,
    onChange?: (o: boolean) => void,
    size?: number,
}

export function MenuButton(props: Props) {
    const {children, size, open, onChange} = props
    const [internalOpen, internalSetOpen] = useState(false)
    const isControlled = open !== undefined
    const realOpen = isControlled ? open : internalOpen

    return (
        <Box sx={SX.box}>
            <Collapse sx={SX.collapse} in={realOpen} orientation={"horizontal"}>
                <Box sx={SX.paper}>{children}</Box>
            </Collapse>
            <MoreIconButton size={size} onClick={() =>  onClick(!realOpen)}/>
        </Box>
    )

    function onClick(value: boolean) {
        if (!isControlled) internalSetOpen(value)
        onChange?.(value)
    }
}
