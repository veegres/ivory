import {PlayArrow, Stop} from "@mui/icons-material"
import {Box, Divider} from "@mui/material"

import {useRouterNodePlatformLogs, useRouterNodePlatformStart, useRouterNodePlatformStop} from "../../../features/node/hook"
import {PlatformConnection} from "../../../features/node/type"
import {UnitBox} from "../../../shared/component/box/UnitBox"
import {SimpleButton} from "../../../shared/component/button/SimpleButton"
import {SxPropsMap} from "../../../shared/helper/type"
import {Logs} from "../logs/Logs"

const SX: SxPropsMap = {
    single: {display: "flex", flexDirection: "column", gap: 0.5, padding: "0px 5px"},
    action: {display: "flex", gap: 0.5},
    head: {display: "flex", justifyContent: "space-between", alignItems: "center", gap: 1},
}

type Props = {
    connection: PlatformConnection,
}

export function MonitorContainerSingle(props: Props) {
    const {connection} = props
    const name = connection.host

    const request = {connection, name, tail: 10, follow: true}
    const logs = useRouterNodePlatformLogs(request)

    const start = useRouterNodePlatformStart(connection)
    const stop = useRouterNodePlatformStop(connection)

    return (
        <Box sx={SX.single}>
            <Divider/>
            <Box sx={SX.head}>
                <Box>{connection.host}</Box>
                <Box sx={SX.action}>
                    <SimpleButton
                        size={"small"}
                        loading={start.isPending}
                        onClick={() => start.mutate({connection, name})}
                    >
                        <PlayArrow fontSize={"small"}/>
                    </SimpleButton>
                    <SimpleButton
                        size={"small"}
                        loading={stop.isPending}
                        onClick={() => stop.mutate({connection, name})}
                    >
                        <Stop fontSize={"small"}/>
                    </SimpleButton>
                </Box>
            </Box>
            <Divider/>
            <UnitBox label={"Logs"} value={logs.data.length} unit={"rows"} fixed={true}/>
            <Logs logs={logs.data} loading={logs.isFetching}/>
        </Box>
    )
}
