import {Box} from "@mui/material"

import {Cluster, Node} from "../../../../features/cluster/api/ClusterType"
import {Feature} from "../../../../features/Feature"
import {ManageAccess} from "../../../../features/management/component/ManageAccess"
import {TabsButton} from "../../../../shared/component/button/TabsButton"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"
import {PgCompactTable} from "../../../../tools/pg_compacttable/component/PgCompactTable"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", justifyContent: "center", gap: 2},
    tabs: {padding: "4px 0px", borderTop: 1, borderBottom: 1, borderColor: "divider"},
}

type Props = {
    node: Node,
    cluster: Cluster,
}

export function NodeMainTools(props: Props) {
    const {node, cluster} = props
    const tab = useStore(s => s.nodeState.toolsTab)
    const {setToolsTab} = useStoreAction
    const tabs = [{label: "pg_compacttable"}]
    return (
        <Box sx={SX.box}>
            <Box sx={SX.tabs}>
                <TabsButton tabs={tabs} tab={tab} setTab={setToolsTab} fullWidth={false}/>
            </Box>
            {tab === 0 && (
                <ManageAccess feature={Feature.ViewToolPgCompactTableList} error={true}>
                    <PgCompactTable node={node} cluster={cluster}/>
                </ManageAccess>
            )}
        </Box>
    )
}