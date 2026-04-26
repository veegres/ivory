import {Box, Divider} from "@mui/material"

import {useRouterClusterOverview} from "../../../../api/cluster/hook"
import {DatabaseType} from "../../../../api/database/type"
import {SxPropsMap} from "../../../../app/type"
import {getQueryConnection} from "../../../../app/utils"
import {useStore, useStoreAction} from "../../../../provider/StoreProvider"
import {AlertCentered} from "../../../view/box/AlertCentered"
import {PageMainBox} from "../../../view/box/PageMainBox"
import {NoDatabaseError} from "../overview/OverviewError"
import {NodeInfo} from "./NodeInfo"
import {NodeMain} from "./NodeMain"

const SX: SxPropsMap = {
    content: {display: "flex", gap: 3},
}

export function Node() {
    const {setNodeBody} = useStoreAction
    const node = useStore(s => s.node)
    const activeCluster = useStore(s => s.activeCluster)
    const activeNodeName = useStore(s => s.activeNode[activeCluster?.cluster.name ?? ""])
    const activeClusterTab = useStore(s => s.activeClusterTab)
    const isClusterOverviewOpen = !!activeCluster && activeClusterTab === 0

    const overview = useRouterClusterOverview(activeCluster?.cluster.name, false)
    const activeNode = overview.data && overview.data.nodes[activeNodeName ?? ""]
    return (
        <PageMainBox withPadding visible={isClusterOverviewOpen}>
            {renderContent()}
        </PageMainBox>
    )

    function renderContent() {
        if (!activeNodeName || !activeCluster) return <AlertCentered text={"Please, select a node to see the information!"}/>
        if (!activeNode) return <AlertCentered text={"There is not enough information about the node!"} severity={"warning"}/>
        if (!activeNode.connection.host || !activeNode.connection.dbPort) return <NoDatabaseError/>

        const database = {type: DatabaseType.POSTGRES, host: activeNode.connection.host, port: activeNode.connection.dbPort}
        const connection = getQueryConnection(activeCluster.cluster, database)

        return (
            <Box sx={SX.content}>
                <NodeInfo node={activeNode} tab={node.body} onTab={setNodeBody} connection={connection}/>
                <Divider orientation={"vertical"} flexItem/>
                <NodeMain tab={node.body} database={database}/>
            </Box>
        )
    }
}
