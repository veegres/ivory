import {Box, LinearProgress, SxProps, Theme} from "@mui/material";
import {SxPropsMap} from "../../../type/general";
import {SxPropsFormatter} from "../../../app/utils";

const SX: SxPropsMap = {
    box: {minHeight: "4px", margin: "5px 0", width: "100%"}
}

type Props = {
    sx?: SxProps<Theme>,
    variant?: "determinate" | "indeterminate" | "buffer" | "query",
    color?: "secondary" | "success" | "inherit" | "warning" | "error" | "primary" | "info",
    line?: boolean,
    isFetching: boolean,
}

export function LinearProgressStateful(props: Props) {
    const {isFetching, sx, variant, color, line} = props
    return (
        <Box sx={SxPropsFormatter.merge(SX.box, sx)}>
            {isFetching ? (
                <LinearProgress variant={variant} color={color}/>
            ) : !line ? null : (
                <LinearProgress color={color} value={0} variant={"determinate"}/>
            )}
        </Box>
    )
}
