import {Box, Collapse, Divider} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../../app/type"

const SX: SxPropsMap = {
    body: {padding: "8px 15px", fontSize: "13px"},
}

type Props = {
    children: ReactNode,
    show: boolean,
    unmountOnExit?: boolean,
}

export function QueryBoxBody(props: Props) {
    const {show, children, unmountOnExit = true} = props

    return (
        <Collapse in={show} unmountOnExit={unmountOnExit} timeout={100}>
            <Divider/>
            <Box sx={SX.body}>
                {children}
            </Box>
        </Collapse>
    )
}
