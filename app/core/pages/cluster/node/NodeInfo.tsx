import {Box} from "@mui/material"

import {Node, NodeConfig} from "../../../../features/cluster/type"
import {SxPropsMap} from "../../../../shared/helper/type"
import {NodeInfoForm} from "./NodeInfoForm"
import {NodeInfoStatus} from "./NodeInfoStatus"

const SX: SxPropsMap = {
    info: {display: "flex", flexWrap: "wrap", columnGap: 1, rowGap: 2, margin: "5px 0", width: "100%"},
}

type Props = {
    node: Node,
    loading?: boolean,
    onUpdate?: (config: NodeConfig) => void,
}

export function NodeInfo(props: Props) {
    const {node, onUpdate, loading} = props

    return (
        <Box sx={SX.info}>
            <NodeInfoStatus role={node.keeper.role}/>
            <NodeInfoForm node={node} onUpdate={onUpdate} loading={loading}/>
        </Box>
    )
}
