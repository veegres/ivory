import {alpha, Box, SxProps, Theme} from "@mui/material"
import {PropsWithChildren} from "react"

import {SxPropsMap} from "../../helper/HelperType"
import {SxPropsFormatter} from "../../helper/HelperUtils"

const SX: SxPropsMap = {
    box: (theme: Theme) => ({
        bgcolor: alpha(theme.palette.primary.light, theme.palette.action.hoverOpacity / 2),
        borderRadius: 2,
        overflow: "hidden",
    }),
}

type Props = PropsWithChildren<{
    sx?: SxProps<Theme>,
}>

export function PaperBlue(props: Props) {
    return <Box sx={SxPropsFormatter.merge(SX.box, props.sx)}>{props.children}</Box>
}
