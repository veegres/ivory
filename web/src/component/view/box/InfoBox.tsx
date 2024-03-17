import {Box, Tooltip, tooltipClasses} from "@mui/material";
import {ReactElement, ReactNode} from "react";
import {SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    box: {
        display: "flex", justifyContent: "center", alignItems: "center", height: "32px",
        fontSize: "0.8125rem", cursor: "pointer", minWidth: "30px", border: 1, borderColor: "divider"
    },
    tooltip: {
        [`& .${tooltipClasses.tooltip}`]: {
            maxWidth: 'none',
        }
    }
}

type Props = {
    tooltip: ReactElement | string,
    children: ReactNode,
    withPadding?: boolean,
    withRadius?: boolean,
}

export function InfoBox(props: Props) {
    const {children, tooltip, withPadding, withRadius} = props
    const padding = withPadding ? "3px 12px" : "3px 5px"
    const borderRadius = withRadius ? "15px" : "4px"

    return (
        <Tooltip PopperProps={{sx: SX.tooltip}} title={tooltip} placement={"top"} arrow={true}>
            <Box sx={{...SX.box, padding, borderRadius}}>
                {children}
            </Box>
        </Tooltip>
    )
}
