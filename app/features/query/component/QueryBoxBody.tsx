import {Box, Collapse, Divider} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../../shared/helper/HelperType"

const SX: SxPropsMap = {
    body: {padding: "8px 15px", fontSize: "13px", backgroundImage: "inherit", backgroundColor: "inherit"},
    inherit: {backgroundImage: "inherit", backgroundColor: "inherit"},
}

type Props = {
    children: ReactNode,
    show: boolean,
    unmountOnExit?: boolean,
}

export function QueryBoxBody(props: Props) {
    const {show, children, unmountOnExit = true} = props

    return (
        <Collapse
            sx={SX.inherit}
            slotProps={{wrapper: {sx: SX.inherit}, wrapperInner: {sx: SX.inherit}}}
            in={show} unmountOnExit={unmountOnExit} timeout={100}
        >
            <Divider/>
            <Box sx={SX.body}>
                {children}
            </Box>
        </Collapse>
    )
}
