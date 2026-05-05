import {Box, Divider} from "@mui/material"

import {useRouterClusterOverview} from "../../../../api/cluster/hook"
import {SxPropsMap} from "../../../../app/type"
import {getQueryConnection, getSshConnection} from "../../../../app/utils"
import {useStore, useStoreAction} from "../../../../provider/StoreProvider"
import {AlertCentered} from "../../../view/box/AlertCentered"
import {PageMainBox} from "../../../view/box/PageMainBox"
import {NodeInfo} from "./NodeInfo"
import {NodeMain} from "./NodeMain"

const SX: SxPropsMap = {
    content: {display: "flex", gap: 3},
}

export function Node() {
    const {setNodeBody} = useStoreAction
    const node = useStore(s => s.nodeState)
    const activeCluster = useStore(s => s.activeCluster)
    const activeClusterName = activeCluster?.name
    const activeNodeName = useStore(s => s.activeNode[activeClusterName ?? ""])

    const overview = useRouterClusterOverview(activeClusterName, false)
    const activeClusterTab = useStore(s => s.activeClusterTab)
    const isClusterOverviewOpen = !!activeCluster && activeClusterTab === 0

    return (
        <PageMainBox withPadding visible={isClusterOverviewOpen}>
            {renderContent()}
        </PageMainBox>
    )

    function renderContent() {
        if (!activeNodeName || !activeClusterName) return <AlertCentered text={"Please, select a node to see the information!"}/>
        const activeNode = overview.data?.nodes[activeNodeName]
        if (!activeNode) return <AlertCentered text={"There is not enough information about the node!"} severity={"warning"}/>
        const {host, dbPort, sshPort, keeperPort} = activeNode.config
        if (!dbPort && !keeperPort && !sshPort) return <AlertCentered text={"Specify at least one port to work with Node"} severity={"warning"}/>

        const sshCon = getSshConnection(activeCluster, host, sshPort)
        const queryCon = getQueryConnection(activeCluster, host, dbPort)

        return (
            <Box sx={SX.content}>
                <NodeInfo node={activeNode} queryCon={queryCon} tab={node.nodeTab} onTab={setNodeBody}/>
                <Divider orientation={"vertical"} flexItem/>
                <NodeMain sshCon={sshCon} queryCon={queryCon} tab={node.nodeTab}/>
            </Box>
        )
    }
}
