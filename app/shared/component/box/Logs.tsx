import {Box, CircularProgress, SxProps, Theme} from "@mui/material"

import {SxPropsMap} from "../../helper/type"
import scroll from "../../style/scroll.module.css"
import {DynamicRowVirtualizer} from "../scrolling/DynamicRowVirtualizer"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1, fontSize: "11px"},
    row: {fontFamily: "monospace", "&:hover": {color: "primary.main"}},
    emptyLine: {textAlign: "center", textTransform: "uppercase"},
    footer: {display: "flex", justifyContent: "space-between", color: "text.secondary"},
    loader: {display: "flex", gap: 0.75, alignItems: "center"},
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
                <Box sx={SX.emptyLine} height={height}>Waiting for logs</Box>
            ) : (
                <Box sx={SX.emptyLine} height={height}>No logs</Box>
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
            <Box sx={SX.footer}>
                <Box sx={SX.loader}>
                    {loading && <CircularProgress sx={SX.loader} size={"9px"} color={"inherit"}/>}
                    <Box>{loading ? "Streaming" : logs.length === 0 ? "None" : "Done"}</Box>
                </Box>
                <Box>{logs.length} rows</Box>
            </Box>
        </Box>
    )
}