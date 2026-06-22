import {Box} from "@mui/material"
import {useState} from "react"

import {Cluster, Node} from "../../../../features/cluster/api/type"
import {Feature} from "../../../../features/feature"
import {ManageAccess} from "../../../../features/management/component/ManageAccess"
import {TabsButton} from "../../../../shared/component/button/TabsButton"
import {SxPropsMap} from "../../../../shared/helper/type"
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
    const [tab, setTab] = useState(0)
    const tabs = [{label: "pg_compacttable"}]
    return (
        <Box sx={SX.box}>
            <Box sx={SX.tabs}>
                <TabsButton tabs={tabs} tab={tab} setTab={setTab} fullWidth={false}/>
            </Box>
            {tab === 0 && (
                <ManageAccess feature={Feature.ViewToolPgCompactTableList} displayError={true}>
                    <PgCompactTable node={node} cluster={cluster}/>
                </ManageAccess>
            )}
        </Box>
    )
}