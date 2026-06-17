import {Box} from "@mui/material"

import {Cluster, Node} from "../../../features/cluster/type"
import {ErrorKeeperRequestMissing} from "../../../shared/component/box/ErrorManual"
import {TitledBox} from "../../../shared/component/box/TitledBox"
import {SxPropsMap} from "../../../shared/helper/type"
import {getKeeperOneRequest} from "../../../shared/helper/utils"
import {KeeperConfig} from "./KeeperConfig"
import {KeeperFailoverButton} from "./KeeperFailoverButton"
import {KeeperReinitButton} from "./KeeperReinitButton"
import {KeeperReloadButton} from "./KeeperReloadButton"
import {KeeperRestartButton} from "./KeeperRestartButton"
import {KeeperScheduleButton} from "./KeeperScheduleButton"
import {KeeperSwitchoverButton} from "./KeeperSwitchoverButton"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", justifyContent: "center", gap: 2},
    actions: {display: "flex", gap: 1, alignItems: "center"},
}

type Props = {
    node: Node,
    cluster: Cluster,
}

export function Keeper(props: Props) {
    const {node, cluster} = props
    const keeperRequest = getKeeperOneRequest(cluster, node.config.host, node.config.keeperPort)

    return (
        <Box sx={SX.box}>
            <TitledBox title={"Actions"} island={true} renderActions={renderActions()}>
                {!keeperRequest && <ErrorKeeperRequestMissing/>}
            </TitledBox>
            <KeeperConfig options={cluster} node={node}/>
        </Box>
    )

    function renderActions() {
        if (!keeperRequest) return
        return (
            <Box sx={SX.actions}>
                <KeeperReloadButton request={keeperRequest} cluster={cluster.name}/>
                <KeeperRestartButton request={keeperRequest} cluster={cluster.name}/>
                <KeeperReinitButton request={keeperRequest} cluster={cluster.name}/>
                <KeeperSwitchoverButton request={keeperRequest} cluster={cluster.name} candidates={cluster.nodes}/>
                <KeeperFailoverButton request={keeperRequest} cluster={cluster.name} role={node.keeper.role}/>
                <KeeperScheduleButton request={keeperRequest} cluster={cluster.name}/>
            </Box>
        )
    }
}