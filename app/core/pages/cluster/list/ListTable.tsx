import {Box, Table, TableCell, TableHead, TableRow} from "@mui/material"
import {useMemo, useState} from "react"

import {ClusterApi} from "../../../../features/cluster/api/router"
import {Cluster} from "../../../../features/cluster/api/type"
import {Feature} from "../../../../features/feature"
import {ManageAccessBox} from "../../../../features/management/component/ManageAccess"
import {KeeperPlugin} from "../../../../features/node/api/type"
import {AlertCentered} from "../../../../shared/component/box/AlertCentered"
import {AddIconButton} from "../../../../shared/component/button/IconButtons"
import {TableBody} from "../../../../shared/component/table/TableBody"
import {TableCellLoader} from "../../../../shared/component/table/TableCellLoader"
import {SxPropsMap} from "../../../../shared/helper/type"
import {SxPropsFormatter} from "../../../../shared/helper/utils"
import {useStore} from "../../../../shared/provider/StoreProvider"
import scroll from "../../../../shared/style/scroll.module.css"
import {Refresher} from "../../../widgets/browser/Refresher"
import {ListDeployCluster} from "./ListDeployCluster"
import {ListDetectCluster} from "./ListDetectCluster"
import {ListRow} from "./ListRow"
import {ListRowNew} from "./ListRowNew"

const SX: SxPropsMap = {
    box: {overflowY: "scroll"},
    table: {"tr:last-child td": {border: 0}, "tr td, th": {padding: "5px 10px"}},
    refresh: {padding: "0px 5px"},
}

type Props = {
    list: Cluster[],
    pending: boolean,
    fetching: boolean,
}

export function ListTable(props: Props) {
    const activeCluster = useStore(s => s.activeCluster)
    const search = useStore(s => s.searchCluster)
    const {list, fetching, pending} = props
    const [showNewElement, setShowNewElement] = useState(false)
    const [editNode, setEditNode] = useState("")

    const rows = useMemo(() => list.filter((c) => c.name.includes(search)), [list, search])

    return (
        <Box sx={SX.box} className={scroll.tiny} maxHeight={activeCluster ? "25vh" : "60vh"}>
            <Table size={"small"} sx={SX.table} stickyHeader>
                <TableHead>
                    <TableRow>
                        <TableCell sx={SxPropsFormatter.style.paper} width={"220px"}>Name</TableCell>
                        <TableCellLoader
                            sx={SxPropsFormatter.style.paper}
                            label={"Nodes"}
                            colSpan={2}
                            loading={fetching && !pending}
                        >
                            <Box sx={SX.refresh}>
                                <Refresher queryKeys={[ClusterApi.list.key(), ClusterApi.overview.key()]}/>
                            </Box>
                            <ListDeployCluster keeper={KeeperPlugin.PATRONI}/>
                            <ListDetectCluster keeper={KeeperPlugin.PATRONI}/>
                            <ManageAccessBox feature={Feature.ManageClusterUpdate}>
                                <AddIconButton
                                    tooltip={"ADD CLUSTER"}
                                    onClick={() => setShowNewElement(true)}
                                    disabled={showNewElement}
                                />
                            </ManageAccessBox>
                        </TableCellLoader>
                    </TableRow>
                </TableHead>
                <TableBody isLoading={pending} cellCount={3} height={32}>
                    <ListRowNew show={showNewElement} close={() => setShowNewElement(false)}/>
                    {renderRemovedRow()}
                    {renderRows()}
                    {renderEmpty()}
                </TableBody>
            </Table>
        </Box>
    )

    function renderRemovedRow() {
        if (!activeCluster) return
        if (rows.some(e => e.name === activeCluster.name)) return
        return (
            <ListRow cluster={activeCluster} editable={false}/>
        )
    }

    function renderRows() {
        return rows.map((cluster) => {
            const editable = cluster.name === editNode
            const toggle = () => setEditNode(editable ? "" : cluster.name)
            return (
                <ListRow key={cluster.name} cluster={cluster} editable={editable} toggle={toggle}/>
            )
        })
    }

    function renderEmpty() {
        if (pending || showNewElement || rows.length || activeCluster) return
        const text = search ? (
            "There are no clusters that match your filter"
        ) : (
            "There are no clusters yet. You can add them manually or use auto detection"
        )
        return (
            <TableRow>
                <TableCell colSpan={3}>
                    <AlertCentered text={text}/>
                </TableCell>
            </TableRow>
        )
    }
}
