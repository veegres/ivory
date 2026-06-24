import {Box, Checkbox, TextField} from "@mui/material"
import {useState} from "react"

import {Logs} from "../../../../shared/component/box/Logs"
import {SxPropsMap} from "../../../../shared/helper/type"
import {useRouterNodePlatformLogs} from "../../api/hook"
import {PlatformConnection} from "../../api/type"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    head: {
        display: "flex", gap: 1, justifyContent: "space-between", padding: "4px 0px", alignItems: "center",
        borderTop: 1, borderBottom: 1, borderColor: "divider", fontSize: "13px",
    },
    options: {display: "flex", gap: 1, alignItems: "center", fontFamily: "monospace"},
    input: {padding: "6px 10px", fontSize: "13px", fontFamily: "monospace"},
    check: {padding: "6px"},
}

type Props = {
    connection: PlatformConnection,
}

export function PlatformLogs(props: Props) {
    const {connection} = props
    const [path, setPath] = useState("")
    const [follow, setFollow] = useState(true)

    const request = {connection, path, tail: 50, follow}
    const logs = useRouterNodePlatformLogs(request)

    return (
        <Box sx={SX.box}>
            <Box sx={SX.head}>
                <TextField
                    size={"small"}
                    variant={"outlined"}
                    placeholder={"Path"}
                    value={path}
                    slotProps={{htmlInput: {sx: SX.input}}}
                    onChange={(e) => setPath(e.target.value)}
                />
                <Box sx={SX.options}>
                    <Box>Follow</Box>
                    <Checkbox
                        size={"small"}
                        checked={follow}
                        slotProps={{root: {sx: SX.check}}}
                        onChange={(e) => setFollow(e.target.checked)}
                    />
                </Box>
            </Box>
            <Logs logs={logs.data} loading={logs.isFetching}/>
        </Box>
    )
}