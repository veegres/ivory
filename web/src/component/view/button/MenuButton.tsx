import {ReactNode, useRef, useState} from "react";
import {Box, Collapse, IconButton, Paper} from "@mui/material";
import {MoreVert} from "@mui/icons-material";
import {SxPropsMap} from "../../../type/general";

const SX: SxPropsMap = {
    box: {position: "relative"},
    collapse: {position: "absolute", top: "50%", transform: "translate(calc(-100% + -2px), -50%)"},
    paper: {display: "flex", gap: "3px", alignItems: "center", padding: "3px 5px", boxShadow: "none", border: 1, borderColor: "divider"},
}

type Props = {
    children: ReactNode,
    size?: "small" | "medium" | "large",
}

export function MenuButton(props: Props) {
    const {children, size} = props
    const [open, setOpen] = useState(false)
    const ref = useRef(null)

    return (
        <Box sx={SX.box}>
            <Collapse sx={SX.collapse} in={open} orientation={"horizontal"}>
                <Paper sx={SX.paper}>
                    {children}
                </Paper>
            </Collapse>
            <IconButton ref={ref} size={size} onClick={() => setOpen(!open)}>
                <MoreVert/>
            </IconButton>
        </Box>
    )

}
