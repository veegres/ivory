import {KeyboardArrowDown} from "@mui/icons-material"
import {Box, Collapse, Typography} from "@mui/material"
import {memo, PropsWithChildren, useState} from "react"

import {SxPropsMap} from "../../../app/type"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column"},
    label: {
        display: "flex", alignItems: "center", cursor: "pointer", userSelect: "none",
        gap: 0.5, padding: "4px 0", "&:hover": {color: "primary.main"},
        "&:hover *": {color: "primary.main"},
    },
    icon: {fontSize: "20px", transition: "transform 0.2s", transform: "rotate(-90deg)"},
    iconOpen: {transform: "rotate(0deg)"},
    text: {
        fontSize: "15px", fontWeight: 600, textTransform: "uppercase", color: "text.secondary",
        fontFamily: "monospace", transition: "color 0.2s", lineHeight: 1,
    },
    content: {marginTop: 2},
}

type Props = {
    label: string,
    defaultOpen?: boolean,
}

export const SubContentBox = memo(function SubContentBox(props: PropsWithChildren<Props>) {
    const {label, children, defaultOpen = false} = props
    const [open, setOpen] = useState(defaultOpen)
    
    return (
        <Box sx={SX.box}>
            <Box sx={SX.label} onClick={() => setOpen(!open)}>
                <KeyboardArrowDown sx={[SX.icon, open && SX.iconOpen]}/>
                <Typography sx={SX.text}>{label}</Typography>
            </Box>
            <Collapse in={open}>
                <Box sx={SX.content}>
                    {children}
                </Box>
            </Collapse>
        </Box>
    )
})
