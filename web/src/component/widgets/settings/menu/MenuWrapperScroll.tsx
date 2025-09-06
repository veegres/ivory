import {ReactNode} from "react";
import {Box, SxProps, Theme} from "@mui/material";
import scroll from "../../../../style/scroll.module.css";
import {SxPropsFormatter} from "../../../../app/utils";

import {SxPropsMap} from "../../../../app/type";

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
