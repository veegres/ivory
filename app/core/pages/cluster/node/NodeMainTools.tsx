import {Box} from "@mui/material"
import {ReactNode} from "react"

import {useRouterClusterOverview} from "../../../../features/cluster/api/ClusterHook"
import {Cluster, Node} from "../../../../features/cluster/api/ClusterType"
import {Feature} from "../../../../features/Feature"
import {ErrorSmart} from "../../../../shared/component/box/ErrorSmart"
import {HeadBox} from "../../../../shared/component/box/HeadBox"
import {TabsButton} from "../../../../shared/component/button/TabsButton"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"
import {PgCompactTable} from "../../../../tools/pg_compacttable/component/PgCompactTable"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", justifyContent: "center", gap: 1},
}

type Tool = {
    label: string,
    feature: Feature,
    render: (node: Node, cluster: Cluster) => ReactNode,
}

// NOTE: add new tools here - each one is hidden automatically unless the
// cluster's plugins support its feature (see Overview.features' doc)
const TOOLS: Tool[] = [
    {
        label: "pg_compacttable",
        feature: Feature.ViewToolPgCompactTableList,
        render: (node, cluster) => <PgCompactTable node={node} cluster={cluster}/>
    },
]

type Props = {
    node: Node,
    cluster: Cluster,
}

export function NodeMainTools(props: Props) {
    const {node, cluster} = props
    const tab = useStore(s => s.nodeState.toolsTab) ?? 0
    const {setToolsTab} = useStoreAction
    const overview = useRouterClusterOverview(cluster.name, false)

    const tools = TOOLS.filter(t => overview.data?.features[t.feature] !== false)
    if (tools.length === 0) return <ErrorSmart error={"No tools are supported by this cluster's plugins"}/>

    return (
        <Box sx={SX.box}>
            <HeadBox>
                <TabsButton tabs={tools} tab={tab} setTab={setToolsTab} fullWidth={false}/>
            </HeadBox>
            {tools[tab]?.render(node, cluster)}
        </Box>
    )
}
