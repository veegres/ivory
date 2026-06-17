import {Box} from "@mui/material"
import {useState} from "react"

import {Cluster, Node} from "../../../features/cluster/type"
import {TabsButton} from "../../../shared/component/button/TabsButton"
import {SxPropsMap} from "../../../shared/helper/type"
import {PgCompactTable} from "../pg_compacttable/PgCompactTable"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", justifyContent: "center", gap: 2},
}

type Props = {
    node: Node,
    cluster: Cluster,
}

export function Tools(props: Props) {
    const {node, cluster} = props
    const [tab, setTab] = useState(0)
    const tabs = [{label: "pg_compacttable"}]
    return (
        <Box sx={SX.box}>
            <TabsButton tabs={tabs} tab={tab} setTab={setTab}/>
            {tab === 0 && <PgCompactTable node={node} cluster={cluster}/>}
        </Box>
    )
}