import {Box, SxProps, Theme} from "@mui/material"

import {SxPropsMap} from "../../helper/HelperType"
import scroll from "../../style/scroll.module.css"
import {DynamicRowVirtualizer} from "../scrolling/DynamicRowVirtualizer"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", fontSize: "11px"},
    row: {fontFamily: "monospace", "&:hover": {color: "primary.main"}},
}

type Props = {
    logs: string[],
    loading?: boolean,
    auto?: boolean,
    sx?: SxProps<Theme>,
    height?: number,
    reconnect?: () => void,
}

export function Logs(props: Props) {
    const {logs, loading = false, auto = true, height = 350, sx, reconnect} = props
    return (
        <Box sx={SX.box}>
            <DynamicRowVirtualizer
                sx={sx}
                sxVirtualRow={SX.row}
                auto={auto}
                className={scroll.small}
                height={height}
                rows={logs}
                reconnect={reconnect}
                title={getTitle()}
            />
        </Box>
    )

    function getTitle() {
        if (logs.length === 0) return loading ? "Waiting for logs" : "No logs"
        return `${loading ? "Streaming" : "Done"} — ${getRows()}`
    }

    function getRows() {
        return logs.length === 1 ? "1 row" : `${logs.length} rows`
    }
}
