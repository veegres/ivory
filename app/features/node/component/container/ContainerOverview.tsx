import {Close, PlayArrow, Replay, Stop} from "@mui/icons-material"
import {Box} from "@mui/material"

import {HeadBox} from "../../../../shared/component/box/HeadBox"
import {Logs} from "../../../../shared/component/box/Logs"
import {AlertButton} from "../../../../shared/component/button/AlertButton"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {useStore} from "../../../../shared/provider/StoreProvider"
import {Feature} from "../../../Feature"
import {ManageAccessBox} from "../../../management/component/ManageAccess"
import {
    useRouterNodePlatformContainerLogs,
    useRouterNodePlatformDown,
    useRouterNodePlatformRestart,
    useRouterNodePlatformStart,
    useRouterNodePlatformStop,
} from "../../api/NodeHook"
import {PlatformVaultConnection} from "../../api/NodeType"
import {ContainerKeeperDeploy} from "./ContainerKeeperDeploy"
import {ContainerOverviewMetrics} from "./ContainerOverviewMetrics"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1, padding: "0px 5px"},
    action: {display: "flex", gap: 0.5},
    logs: {
        display: "flex", flexDirection: "column", gap: 0.5, padding: "5px",
        border: 1, borderRadius: 1, borderColor: "divider",
    }
}

type Props = {
    connection: PlatformVaultConnection,
    name: string,
}

export function ContainerOverview(props: Props) {
    const {connection, name} = props
    const activeCluster = useStore(s => s.activeCluster)

    const request = {connection, path: name, tail: 50, follow: true}
    const logs = useRouterNodePlatformContainerLogs(request)

    const start = useRouterNodePlatformStart(connection)
    const stop = useRouterNodePlatformStop(connection)
    const restart = useRouterNodePlatformRestart(connection)
    const down = useRouterNodePlatformDown(connection)

    return (
        <Box sx={SX.box}>
            <HeadBox title={name}>
                <ManageAccessBox sx={SX.action} feature={Feature.ManageNodePlatformContainer}>
                    {activeCluster && (
                        <ContainerKeeperDeploy
                            connection={connection}
                            plugin={activeCluster.plugins.keeper}
                            cluster={activeCluster.name}
                            node={name}
                            databaseId={activeCluster.vaults.databaseId}
                            sshKeyId={activeCluster.vaults.sshKeyId}
                        />
                    )}
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
                </ManageAccessBox>
            </HeadBox>
            <ContainerOverviewMetrics connection={connection} name={name}/>
            <Box sx={SX.logs}>
                <Logs logs={logs.data} loading={logs.isFetching} reconnect={logs.reconnect}/>
            </Box>
        </Box>
    )
}
