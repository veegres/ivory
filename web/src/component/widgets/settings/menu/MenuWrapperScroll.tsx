import {Box, SxProps, Theme} from "@mui/material";
import {ReactNode} from "react";

import {SxPropsMap} from "../../../../app/type";
import {SxPropsFormatter} from "../../../../app/utils";
import scroll from "../../../../style/scroll.module.css";

const SX: SxPropsMap = {
    box: {height: "100%", overflowY: "auto", padding: "0 10px"},
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
