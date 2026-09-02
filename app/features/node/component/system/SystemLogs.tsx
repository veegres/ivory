import {Box, TextField, ToggleButton, Tooltip} from "@mui/material"
import {useState} from "react"

import {Logs} from "../../../../shared/component/box/Logs"
import {NoBox} from "../../../../shared/component/box/NoBox"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {useDebounce} from "../../../../shared/hook/Debounce"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"
import {useRouterNodeSystemLogs} from "../../api/NodeHook"
import {PlatformVaultConnection} from "../../api/NodeType"

const DEFAULT_TAIL = 50

const SX: SxPropsMap = {
    head: {display: "flex", gap: 1, alignItems: "center"},
    path: {flexGrow: 1},
    tail: {width: "90px"},
    follow: {padding: "0px 10px", border: 1, borderColor: "divider", borderRadius: 1, lineHeight: 1},
    input: {fontFamily: "monospace", fontSize: "13px"},
}

type Props = {
    connection: PlatformVaultConnection,
}

export function SystemLogs(props: Props) {
    const {connection} = props
    const path = useStore(s => s.nodeState.systemLogsPath)
    const {setSystemLogsPath} = useStoreAction
    const [tail, setTail] = useState("")
    const [follow, setFollow] = useState(true)
    const debouncePath = useDebounce(path)
    const debounceTail = useDebounce(tail)

    const request = {connection, path: debouncePath, tail: getTail(debounceTail), follow}
    const logs = useRouterNodeSystemLogs(request, debouncePath !== "")

    return (
        <>
            <Box sx={SX.head}>
                <TextField
                    sx={SX.path}
                    size={"small"}
                    color={"secondary"}
                    placeholder={"Path"}
                    value={path}
                    slotProps={{htmlInput: {sx: SX.input}}}
                    onChange={(e) => setSystemLogsPath(e.target.value)}
                />
                <Tooltip title={"Number of last rows"} placement={"top"}>
                    <TextField
                        sx={SX.tail}
                        size={"small"}
                        color={"secondary"}
                        placeholder={DEFAULT_TAIL.toString()}
                        value={tail}
                        slotProps={{htmlInput: {sx: SX.input}}}
                        onChange={(e) => setTail(e.target.value.replace(/\D/g, ""))}
                    />
                </Tooltip>
                <Tooltip title={"Keep streaming new rows"} placement={"top"}>
                    <ToggleButton
                        sx={SX.follow}
                        size={"small"}
                        value={"follow"}
                        color={"secondary"}
                        selected={follow}
                        onChange={() => setFollow(!follow)}
                    >
                        Follow
                    </ToggleButton>
                </Tooltip>
            </Box>
            {path === "" ? (
                <NoBox text={"Enter a file path to view logs"}/>
            ) : (
                <Logs logs={logs.data} loading={logs.isFetching} reconnect={logs.reconnect}/>
            )}
        </>
    )

    function getTail(value: string) {
        const rows = parseInt(value)
        return isNaN(rows) || rows <= 0 ? DEFAULT_TAIL : rows
    }
}
