import {ReactNode, useRef, useState} from "react";
import {Box, IconButton, Paper, Popper} from "@mui/material";
import {MoreVert} from "@mui/icons-material";
import {SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    paper: {display: "flex", gap: "3px", alignItems: "center", padding: "3px 5px"},
}

type Props = {
    children: ReactNode,
    size?: "small" | "medium" | "large",
}

export function MenuButton(props: Props) {
    const {children, size} = props
    const [open, setOpen] = useState(false)
    const ref = useRef(null)
    const preventOverFlow = {name: 'preventOverflow', options: {mainAxis: false}}

    return (
        <Box>
            <Popper anchorEl={ref.current} open={open} placement={"left"} modifiers={[preventOverFlow]}>
                <Paper sx={SX.paper} variant={"outlined"}>
                    {children}
                </Paper>
            </Popper>
            <IconButton ref={ref} size={size} onClick={() => setOpen(!open)}>
                <MoreVert/>
            </IconButton>
        </Box>
    )

}
