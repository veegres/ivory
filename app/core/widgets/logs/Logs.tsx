import {Box, SxProps, Theme} from "@mui/material"

import {LinearProgressStateful} from "../../../shared/component/progress/LinearProgressStateful"
import {DynamicRowVirtualizer} from "../../../shared/component/scrolling/DynamicRowVirtualizer"
import {SxPropsMap} from "../../../shared/helper/type"
import scroll from "../../../shared/style/scroll.module.css"

const SX: SxPropsMap = {
    box: {fontSize: "12px"},
    row: {"&:hover": {color: "primary.main"}},
    emptyLine: {textAlign: "center", textTransform: "uppercase"},
    loader: {margin: "10px 0 5px"},
}

type Props = {
    logs: string[],
    loading?: boolean,
    auto?: boolean,
    sx?: SxProps<Theme>,
    height?: number,
}

export function Logs(props: Props) {
    const {logs, loading = false, auto = true, height = 350, sx} = props
    return (
        <Box sx={SX.box}>
            {logs.length === 0 ? loading ? (
                <Box sx={SX.emptyLine}>Waiting for logs</Box>
            ) : (
                <Box sx={SX.emptyLine}>No logs</Box>
            ) : (
                <DynamicRowVirtualizer
                    sx={sx}
                    sxVirtualRow={SX.row}
                    auto={auto}
                    className={scroll.small}
                    height={height}
                    rows={logs}
                />
            )}
            <LinearProgressStateful sx={SX.loader} loading={loading} color={"inherit"} line/>
        </Box>
    )
}