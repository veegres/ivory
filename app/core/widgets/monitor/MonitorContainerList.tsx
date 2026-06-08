import {Box} from "@mui/material"

import {useRouterNodePlatformList} from "../../../features/node/hook"
import {PlatformConnection} from "../../../features/node/type"
import {ErrorSmart} from "../../../shared/component/box/ErrorSmart"
import {SxPropsMap} from "../../../shared/helper/type"
import scroll from "../../../shared/style/scroll.module.css"
import {MonitorLoading} from "./MonitorLoading"

const SX: SxPropsMap = {
    list: {
        display: "flex", flexDirection: "column", gap: 0.5, padding: "10px 10px 5px",
        overflowX: "scroll", overflowY: "auto", maxHeight: "400px", minHeight: "100px",
        backgroundColor: "background.default", color: "text.secondary", borderRadius: 2,
        border: 0.5, borderColor: "divider",
    },
    pre: {margin: 0, fontFamily: "'Fira Code', 'Courier New', monospace", fontSize: "12px", whiteSpace: "pre"}
}

type Props = {
    connection: PlatformConnection,
}

export function MonitorContainerList(props: Props) {
    const {connection} = props
    const list = useRouterNodePlatformList(connection)

    if (list.isError) return <ErrorSmart error={list.error}/>
    if (list.isPending) return <MonitorLoading count={5}/>

    return (
        <Box sx={SX.list} className={scroll.tiny}>
            <Box sx={SX.pre}>
                {list.data?.join("\n")}
            </Box>
        </Box>
    )
}
