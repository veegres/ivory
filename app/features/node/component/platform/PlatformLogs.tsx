import {Box, Checkbox, TextField} from "@mui/material"
import {useState} from "react"

import {Logs} from "../../../../shared/component/box/Logs"
import {NoBox} from "../../../../shared/component/box/NoBox"
import {SxPropsMap} from "../../../../shared/helper/type"
import {useDebounce} from "../../../../shared/hook/Debounce"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"
import {useRouterNodePlatformLogs} from "../../api/hook"
import {PlatformVaultConnection} from "../../api/type"

const SX: SxPropsMap = {
    head: {
        display: "flex", gap: 1, justifyContent: "space-between", padding: "4px 0px", alignItems: "center",
        borderTop: 1, borderBottom: 1, borderColor: "divider", fontSize: "13px",
    },
    options: {display: "flex", gap: 1, alignItems: "center", fontFamily: "monospace"},
    input: {padding: "6px 10px", fontSize: "14px", fontFamily: "monospace"},
    check: {padding: "6px"},
}

type Props = {
    connection: PlatformVaultConnection,
}

export function PlatformLogs(props: Props) {
    const {connection} = props
    const path = useStore(s => s.nodeState.platformLogsPath)
    const {setPlatformLogsPath} = useStoreAction
    const debouncePath = useDebounce(path)
    const [follow, setFollow] = useState(true)

    const request = {connection, path: debouncePath, tail: 50, follow}
    const logs = useRouterNodePlatformLogs(request, debouncePath !== "")

    return (
        <>
            <Box sx={SX.head}>
                <TextField
                    size={"small"}
                    variant={"outlined"}
                    placeholder={"Path"}
                    value={path}
                    slotProps={{htmlInput: {sx: SX.input}}}
                    onChange={(e) => setPlatformLogsPath(e.target.value)}
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
            {path === "" ? (
                <NoBox text={"enter path to see logs"}/>
            ) : (
                <Logs logs={logs.data} loading={logs.isFetching} reconnect={logs.reconnect}/>
            )}
        </>
    )
}