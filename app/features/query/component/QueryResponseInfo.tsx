import {Box, Tooltip} from "@mui/material"

import {SxPropsMap} from "../../../shared/helper/type"
import {DbOptions} from "../api/type"

const SX: SxPropsMap = {
    box: {
        display: "flex", justifyContent: "space-between", alignItems: "center", gap: 2, width: "100%",
        color: "text.secondary", cursor: "pointer", whiteSpace: "nowrap",
    },
    right: {display: "flex", alignItems: "center", gap: 1},
    wrap: {whiteSpace: "wrap", wordBreak: "break-all"}
}

type Props = {
    url: string,
    time: {start?: number, end?: number},
    options?: DbOptions,
}

export function QueryResponseInfo(props: Props) {
    const {url, options, time} = props
    return (
        <Box sx={SX.box}>
            <Tooltip title={"SENT TO"} placement={"right"} arrow={true}>
                <Box sx={SX.wrap}>[ {url} ]</Box>
            </Tooltip>
            <Box sx={SX.right}>
                <Tooltip title={"LIMIT"} placement={"top"}>
                    <Box>[ {options?.limit ?? "-"} ]</Box>
                </Tooltip>
                <Tooltip title={"DURATION"} placement={"top"}>
                    <Box>[ {time.end && time.start ? `${(time.end - time.start) / 1000}s` : "-"} ]</Box>
                </Tooltip>
            </Box>
        </Box>
    )
}
