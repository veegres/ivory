import {Close, PlayArrow, Replay, Stop} from "@mui/icons-material"
import {Box, Divider} from "@mui/material"

import {
    useRouterNodePlatformDown,
    useRouterNodePlatformLogs,
    useRouterNodePlatformRestart,
    useRouterNodePlatformStart,
    useRouterNodePlatformStop,
} from "../../../features/node/hook"
import {PlatformConnection} from "../../../features/node/type"
import {AlertButton} from "../../../shared/component/button/AlertButton"
import {SxPropsMap} from "../../../shared/helper/type"
import {useStore} from "../../../shared/provider/StoreProvider"
import {Logs} from "../logs/Logs"
import {MonitorContainerDeploy} from "./MonitorContainerDeploy"

const SX: SxPropsMap = {
    single: {display: "flex", flexDirection: "column", gap: 0.5, padding: "0px 5px"},
    action: {display: "flex", gap: 0.5},
    head: {display: "flex", justifyContent: "space-between", alignItems: "center", gap: 1},
    name: {
        fontFamily: "monospace", fontSize: "13px", border: 1, borderRadius: 1,
        borderColor: "divider", padding: "4px 10px",
    },
}

type Props = {
    connection: PlatformConnection,
}

export function MonitorContainerSingle(props: Props) {
    const {connection} = props
    const name = connection.host

    const activeCluster = useStore(s => s.activeCluster)

    const request = {connection, name, tail: 50, follow: true}
    const logs = useRouterNodePlatformLogs(request)

    const start = useRouterNodePlatformStart(connection)
    const stop = useRouterNodePlatformStop(connection)
    const restart = useRouterNodePlatformRestart(connection)
    const down = useRouterNodePlatformDown(connection)

    return (
        <Box sx={SX.single}>
            <Divider/>
            <Box sx={SX.head}>
                <Box sx={SX.name}>{connection.host}</Box>
                <Box sx={SX.action}>
                    {activeCluster && <MonitorContainerDeploy connection={connection} cluster={activeCluster}/>}
                    <AlertButton
                        color={"inherit"}
                        variant={"outlined"}
                        label={<PlayArrow fontSize={"small"}/>}
                        title={"Start container"}
                        tooltip={"START CONTAINER"}
                        loading={start.isPending}
                        onClick={() => start.mutate({connection, name})}
                        description={"This will start the container."}
                    />
                    <AlertButton
                        color={"inherit"}
                        variant={"outlined"}
                        label={<Replay fontSize={"small"}/>}
                        title={"Restart container"}
                        tooltip={"RESTART CONTAINER"}
                        loading={restart.isPending}
                        onClick={() => restart.mutate({connection, name})}
                        description={"This will restart the container."}
                    />
                    <AlertButton
                        color={"inherit"}
                        variant={"outlined"}
                        label={<Stop fontSize={"small"}/>}
                        title={"Stop container"}
                        tooltip={"STOP CONTAINER"}
                        loading={stop.isPending}
                        onClick={() => stop.mutate({connection, name})}
                        description={"This will stop the container."}
                    />
                    <AlertButton
                        color={"inherit"}
                        variant={"outlined"}
                        label={<Close fontSize={"small"}/>}
                        title={"Remove container"}
                        tooltip={"REMOVE CONTAINER"}
                        loading={down.isPending}
                        onClick={() => down.mutate({connection, name})}
                        description={"This will remove container service. But you need first to stop it."}
                    />
                </Box>
            </Box>
            <Divider/>
            <Logs logs={logs.data} loading={logs.isFetching}/>
        </Box>
    )
}
