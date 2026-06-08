import {Close, PlayArrow, Rocket, Stop} from "@mui/icons-material"
import {Box, Divider} from "@mui/material"

import {
    useRouterNodePlatformDown,
    useRouterNodePlatformLogs,
    useRouterNodePlatformStart,
    useRouterNodePlatformStop,
} from "../../../features/node/hook"
import {PlatformConnection} from "../../../features/node/type"
import {UnitBox} from "../../../shared/component/box/UnitBox"
import {SimpleButton} from "../../../shared/component/button/SimpleButton"
import {SxPropsMap} from "../../../shared/helper/type"
import {Logs} from "../logs/Logs"

const SX: SxPropsMap = {
    single: {display: "flex", flexDirection: "column", gap: 0.5, padding: "0px 5px"},
    action: {display: "flex", gap: 0.5},
    head: {display: "flex", justifyContent: "space-between", alignItems: "center", gap: 1},
    name: {
        fontFamily: "monospace", fontSize: "13px", border: 1, borderRadius: 1,
        borderColor: "divider", padding: "5px 10px",
    },
}

type Props = {
    connection: PlatformConnection,
}

export function MonitorContainerSingle(props: Props) {
    const {connection} = props
    const name = connection.host

    const request = {connection, name, tail: 50, follow: true}
    const logs = useRouterNodePlatformLogs(request)

    const start = useRouterNodePlatformStart(connection)
    const stop = useRouterNodePlatformStop(connection)
    const down = useRouterNodePlatformDown(connection)

    return (
        <Box sx={SX.single}>
            <Divider/>
            <Box sx={SX.head}>
                <Box sx={SX.name}>{connection.host}</Box>
                <Box sx={SX.action}>
                    <SimpleButton
                        size={"small"}
                        loading={false}
                        onClick={() => {}}
                        disabled={true}
                    >
                        <Rocket fontSize={"small"}/>
                    </SimpleButton>
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
                    <SimpleButton
                        size={"small"}
                        loading={down.isPending}
                        onClick={() => down.mutate({connection, name})}
                    >
                        <Close fontSize={"small"}/>
                    </SimpleButton>
                </Box>
            </Box>
            <Divider/>
            <UnitBox label={"Logs"} value={logs.data.length} unit={"rows"} fixed={true}/>
            <Logs logs={logs.data} loading={logs.isFetching}/>
        </Box>
    )
}
