import {Box} from "@mui/material"

import {useRouterClusterUpdate} from "../../../../features/cluster/api/ClusterHook"
import {useRouterClusterOverview} from "../../../../features/cluster/api/ClusterHook"
import {NodeConfig} from "../../../../features/cluster/api/ClusterType"
import {AlertCentered} from "../../../../shared/component/box/AlertCentered"
import {PageMainBox} from "../../../../shared/component/box/PageMainBox"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {getDomain} from "../../../../shared/helper/HelperUtils"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"
import {NodeHead} from "./NodeHead"
import {NodeMain} from "./NodeMain"

const SX: SxPropsMap = {
    content: {display: "flex", flexDirection: "column", backgroundImage: "inherit", backgroundColor: "inherit"},
}

export function Node() {
    const activeCluster = useStore(s => s.activeCluster)
    const activeClusterName = activeCluster?.name
    const activeNodeName = useStore(s => s.activeNode[activeClusterName ?? ""])

    const overview = useRouterClusterOverview(activeClusterName, false)
    const updateCluster = useRouterClusterUpdate(activeCluster?.name!)

    return (
        <PageMainBox withPadding visible={!!activeCluster}>
            {renderContent()}
        </PageMainBox>
    )

    function renderContent() {
        if (!activeNodeName || !activeClusterName) return <AlertCentered text={"Please, select a node to see the information!"}/>
        const activeNode = overview.data?.nodes[activeNodeName]
        if (!activeNode) return <AlertCentered text={"There is not enough information about the node!"} severity={"warning"}/>
        const {dbPort, sshPort, keeperPort} = activeNode.config
        if (!dbPort && !keeperPort && !sshPort) return <AlertCentered text={"Specify at least one port to work with Node"} severity={"warning"}/>

        return (
            <Box sx={SX.content}>
                <NodeHead node={activeNode} loading={updateCluster.isPending} onUpdate={(c) => handleUpdateNode(c, activeNode.config.host)}/>
                <NodeMain cluster={activeCluster} node={activeNode}/>
            </Box>
        )
    }

    function handleUpdateNode(config: NodeConfig, host: string) {
        if (!activeCluster) return
        const nodes = activeCluster.nodes.map(n => n.host === host ? config : n)
        updateCluster.mutate({...activeCluster, nodes})
        // NOTE: this should be done only on success, but it is ok for now, can be improved later
        useStoreAction.setNode(getDomain(config, true))
    }
}
