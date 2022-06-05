import {Box, LinearProgress} from "@mui/material";
import {SxProps} from "@mui/system";

const SX = {
    box: { minHeight: "4px" }
}

type Props = {
    sx?: SxProps,
    variant?: "determinate" | "indeterminate" | "buffer" | "query",
    color?: "secondary" | "success" | "inherit" | "warning" | "error" | "primary" | "info",
    line?: boolean,
    isFetching: boolean,
}

export function LinearProgressStateful(props: Props) {
    const { isFetching, sx, variant, color, line } = props
    return (
        <Box sx={{...SX.box, ...sx}}>
            {isFetching ? (
                <LinearProgress variant={variant} color={color} />
            ) : !line ? null : (
                <LinearProgress color={color} value={0} variant={"determinate"} />
            )}
        </Box>
    )
}
