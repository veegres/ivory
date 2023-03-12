import {Box, Tooltip} from "@mui/material";
import {ReactElement, ReactNode} from "react";
import {useTheme} from "../../provider/ThemeProvider";
import {SxPropsMap} from "../../type/common";

const SX: SxPropsMap = {
    box: { display: "flex", justifyContent: "center", alignItems: "center", height: "32px", fontSize: "0.8125rem", cursor: "pointer", minWidth: "30px" },
}

type Props = {
    tooltip: ReactElement | string,
    children: ReactNode,
    withPadding?: boolean,
    withRadius?: boolean,
}

export function InfoBox(props: Props) {
    const { children, tooltip, withPadding, withRadius } = props
    const { info } = useTheme()
    const padding = withPadding ? "3px 12px" : "3px 5px"
    const borderRadius = withRadius ? "15px" : "4px"

    return (
        <Tooltip title={tooltip} placement={"top"} arrow>
            <Box sx={{ ...SX.box, padding, borderRadius, border: `1px solid ${info?.palette.divider}` }}>
                {children}
            </Box>
        </Tooltip>
    )
}
