import {Box} from "@mui/material"

import {Cluster, Node} from "../../../../features/cluster/api/ClusterType"
import {Feature} from "../../../../features/Feature"
import {ManageAccess} from "../../../../features/management/component/ManageAccess"
import {HeadBox} from "../../../../shared/component/box/HeadBox"
import {TabsButton} from "../../../../shared/component/button/TabsButton"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"
import {PgCompactTable} from "../../../../tools/pg_compacttable/component/PgCompactTable"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", justifyContent: "center", gap: 1},
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
            <HeadBox>
                <TabsButton tabs={tabs} tab={tab} setTab={setToolsTab} fullWidth={false}/>
            </HeadBox>
            {tab === 0 && (
                <ManageAccess feature={Feature.ViewToolPgCompactTableList} error={true}>
                    <PgCompactTable node={node} cluster={cluster}/>
                </ManageAccess>
            )}
        </Box>
    )
}