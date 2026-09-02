import {Close, PlayArrow, Replay, Stop} from "@mui/icons-material"

import {HeadBox} from "../../../../shared/component/box/HeadBox"
import {AlertButton} from "../../../../shared/component/button/AlertButton"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {useStore} from "../../../../shared/provider/StoreProvider"
import {Feature} from "../../../Feature"
import {ManageAccessBox} from "../../../management/component/ManageAccess"
import {
    useRouterNodePlatformDown,
    useRouterNodePlatformRestart,
    useRouterNodePlatformStart,
    useRouterNodePlatformStop,
} from "../../api/NodeHook"
import {PlatformVaultConnection} from "../../api/NodeType"
import {ContainerKeeperDeploy} from "./ContainerKeeperDeploy"

const SX: SxPropsMap = {
    action: {display: "flex", gap: 0.5},
}

type Props = {
    connection: PlatformVaultConnection,
    name: string,
}

export function ContainerHead(props: Props) {
    const {connection, name} = props
    const activeCluster = useStore(s => s.activeCluster)

    const start = useRouterNodePlatformStart(connection)
    const stop = useRouterNodePlatformStop(connection)
    const restart = useRouterNodePlatformRestart(connection)
    const down = useRouterNodePlatformDown(connection)

    return (
        <HeadBox title={name}>
            <ManageAccessBox sx={SX.action} feature={Feature.ManageNodePlatformContainer}>
                {activeCluster && (
                    <ContainerKeeperDeploy
                        connection={connection}
                        plugin={activeCluster.plugins.keeper}
                        cluster={activeCluster.name}
                        node={name}
                        keeperId={activeCluster.vaults.keeperId}
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
    )
}
