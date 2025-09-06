import {ReactNode, useState} from "react";
import {Box, Collapse, Paper} from "@mui/material";
import {SxPropsMap} from "../../../api/management/type";
import {MoreIconButton} from "./IconButtons";

const SX: SxPropsMap = {
    box: {position: "relative"},
    collapse: {position: "absolute", top: "50%", transform: "translate(calc(-100% + -2px), -50%)"},
    paper: {display: "flex", gap: "3px", alignItems: "center", padding: "3px 5px", boxShadow: "none", border: 1, borderColor: "divider"},
}

type Props = {
    children: ReactNode,
    size?: number,
}

export function MenuButton(props: Props) {
    const {children, size} = props
    const [open, setOpen] = useState(false)

    return (
        <Box sx={SX.box}>
            <Collapse sx={SX.collapse} in={open} orientation={"horizontal"}>
                <Paper sx={SX.paper}>
                    {children}
                </Paper>
            </Collapse>
            <MoreIconButton size={size} onClick={() => setOpen(!open)}/>
        </Box>
    )

}
