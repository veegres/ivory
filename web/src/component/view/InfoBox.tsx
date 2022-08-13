import {Box, Tooltip} from "@mui/material";
import {ReactElement, ReactNode} from "react";
import {useTheme} from "../../provider/ThemeProvider";

const SX = {
    box: { display: "flex", alignItems: "center", borderRadius: "4px", height: "32px", fontSize: "0.8125rem" },
}

type Props = {
    tooltip: ReactElement | string,
    children: ReactNode,
    withPadding?: boolean,
}

export function InfoBox(props: Props) {
    const { children, tooltip, withPadding } = props
    const { info } = useTheme()
    const padding = withPadding ? "3px 12px" : "3px 5px"

    return (
        <Tooltip title={tooltip} placement={"top"}>
            <Box sx={{ ...SX.box, padding, border: `1px solid ${info?.palette.divider}` }}>
                {children}
            </Box>
        </Tooltip>
    )
}
