import {Box, LinearProgress, SxProps, Theme} from "@mui/material";

import {SxPropsMap} from "../../../app/type";
import {SxPropsFormatter} from "../../../app/utils";

const SX: SxPropsMap = {
    box: {minHeight: "4px", margin: "5px 0", width: "100%"},
}

type Props = {
    sx?: SxProps<Theme>,
    variant?: "determinate" | "indeterminate" | "buffer" | "query",
    color?: "secondary" | "success" | "inherit" | "warning" | "error" | "primary" | "info",
    line?: boolean,
    loading: boolean,
}

export function LinearProgressStateful(props: Props) {
    const {loading, sx, variant, color, line} = props
    return (
        <Box sx={SxPropsFormatter.merge(SX.box, sx)}>
            {loading ? (
                <LinearProgress variant={variant} color={color}/>
            ) : !line ? null : (
                <LinearProgress color={color} value={0} variant={"determinate"}/>
            )}
        </Box>
    )
}
