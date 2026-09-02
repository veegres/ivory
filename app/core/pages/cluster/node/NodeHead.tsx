import {Box} from "@mui/material"

import {Node, NodeConfig} from "../../../../features/cluster/api/ClusterType"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {NodeHeadForm} from "./NodeHeadForm"
import {NodeHeadRole} from "./NodeHeadRole"

const SX: SxPropsMap = {
    info: {display: "flex", flexWrap: "wrap", columnGap: 1, rowGap: 2, margin: "5px 0", width: "100%"},
}

type Props = {
    node: Node,
    loading?: boolean,
    onUpdate?: (config: NodeConfig) => void,
}

export function NodeHead(props: Props) {
    const {node, onUpdate, loading} = props

    return (
        <Box sx={SX.info}>
            <NodeHeadRole role={node.keeper.role}/>
            <NodeHeadForm node={node} onUpdate={onUpdate} loading={loading}/>
        </Box>
    )
}
