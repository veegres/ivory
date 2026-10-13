import {Box} from "@mui/material"

import {Cluster, NodeOverview} from "../../../../features/cluster/api/ClusterType"
import {KeeperLag} from "../../../../features/node/api/NodeType"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {getInitialNode, getKeeperCandidates, getNodeConfig} from "../../../../shared/helper/HelperUtils"
import {useStore} from "../../../../shared/provider/StoreProvider"
import {OverviewNodesRow} from "./OverviewNodesRow"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 2},
    table: {display: "flex", flexDirection: "column"},
}

type Props = {
    cluster: Cluster,
    nodes?: NodeOverview,
    lag?: KeeperLag,
}

export function OverviewNodes(props: Props) {
    const {cluster, lag} = props
    const nodes = props.nodes ?? cluster.nodesOverview
    const activeNode = useStore(s => s.activeNode[cluster.name])
    const candidates = getKeeperCandidates(nodes)

    return (
        <Box sx={SX.box}>
            <Box sx={SX.table}>
                {renderCheckedNode()}
                {Object.entries(nodes ?? {}).map(([key, element]) => (
                    <OverviewNodesRow
                        key={key}
                        nodeKey={key}
                        checked={key === activeNode}
                        node={element}
                        cluster={cluster}
                        candidates={candidates}
                        lag={lag}
                    />
                ))}
            </Box>
        </Box>
    )

    function renderCheckedNode() {
        if (!activeNode) return
        if (nodes?.[activeNode]) return
        return (
            <OverviewNodesRow
                nodeKey={activeNode}
                node={getInitialNode(getNodeConfig(activeNode))}
                checked={true}
                cluster={cluster}
                candidates={candidates}
                lag={lag}
                error={true}
            />
        )
    }
}
