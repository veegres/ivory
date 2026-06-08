import {Box} from "@mui/material"

import {Node, NodeConfig} from "../../../../features/cluster/type"
import {SxPropsMap} from "../../../../shared/helper/type"
import {NodeInfoStatus} from "./NodeInfoStatus"
import {NodeInfoTable} from "./NodeInfoTable"

const SX: SxPropsMap = {
    info: {display: "flex", gap: 1, margin: "5px 0"},
}

type Props = {
    node: Node,
    onUpdate?: (config: NodeConfig) => void,
}

export function NodeInfo(props: Props) {
    const {node, onUpdate} = props

    return (
        <Box sx={SX.info}>
            <NodeInfoStatus role={node.keeper.role}/>
            <NodeInfoTable node={node} onUpdate={onUpdate}/>
        </Box>
    )
}
