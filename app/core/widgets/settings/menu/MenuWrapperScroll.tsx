import {Box, SxProps, Theme} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../../../shared/helper/type"
import {SxPropsFormatter} from "../../../../shared/helper/utils"
import scroll from "../../../../shared/style/scroll.module.css"

const SX: SxPropsMap = {
    box: {height: "100%", overflowY: "scroll", padding: "0px 5px 0px 0px"},
}

type Props = {
    sx?: SxProps<Theme>,
    children: ReactNode,
}

export function MenuWrapperScroll(props: Props) {
    const {sx, children} = props
    return (
        <Box sx={SxPropsFormatter.merge(SX.box, sx)} className={scroll.small}>
            {children}
        </Box>
    )
}
