import {Box, CircularProgress, SxProps, Theme} from "@mui/material"

import {DynamicRowVirtualizer} from "../../../shared/component/scrolling/DynamicRowVirtualizer"
import {SxPropsMap} from "../../../shared/helper/type"
import scroll from "../../../shared/style/scroll.module.css"

const SX: SxPropsMap = {
    box: {fontSize: "11px"},
    row: {fontFamily: "monospace", "&:hover": {color: "primary.main"}},
    emptyLine: {textAlign: "center", textTransform: "uppercase"},
    footer: {display: "flex", justifyContent: "space-between", padding: "5px 0px", color: "text.secondary"},
    loader: {display: "flex", gap: 1, alignItems: "center"},
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
            <Box sx={SX.footer}>
                <Box sx={SX.loader}>
                    {loading && <CircularProgress sx={SX.loader} size={"10px"} color={"inherit"}/>}
                    {loading ? "Loading" : "Finished"}
                </Box>
                <Box>{logs.length} rows</Box>
            </Box>
        </Box>
    )
}