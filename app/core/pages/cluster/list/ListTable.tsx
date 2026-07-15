import {Box, Table, TableCell, TableHead, TableRow} from "@mui/material"
import {useMemo, useState} from "react"

import {Cluster} from "../../../../features/cluster/api/ClusterType"
import {ClusterDeploy} from "../../../../features/cluster/component/ClusterDeploy"
import {ClusterDetect} from "../../../../features/cluster/component/ClusterDetect"
import {AlertCentered} from "../../../../shared/component/box/AlertCentered"
import {TableBody} from "../../../../shared/component/table/TableBody"
import {TableCellLoader} from "../../../../shared/component/table/TableCellLoader"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {KeeperPluginOptions, SxPropsFormatter} from "../../../../shared/helper/HelperUtils"
import {useStore} from "../../../../shared/provider/StoreProvider"
import scroll from "../../../../shared/style/scroll.module.css"
import {ListClusterAdd} from "./ListClusterAdd"
import {ListEmptyInfo} from "./ListEmptyInfo"
import {ListRow} from "./ListRow"
import {ListRowNew} from "./ListRowNew"
import {ListTableRefresher} from "./ListTableRefresher"

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
    const keeper = useStore(s => s.activeClusterKeeperPlugin)
    const activeCluster = useStore(s => s.activeCluster)
    const search = useStore(s => s.searchCluster)
    const {list, fetching, pending} = props
    const [showNewElement, setShowNewElement] = useState(false)
    const [editNode, setEditNode] = useState("")

    const rows = useMemo(() => list.filter((c) => c.name.includes(search)), [list, search])

    return (
        <Box sx={[SX.box, {maxHeight: activeCluster ? "25vh" : "60vh"}]} className={scroll.tiny}>
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
                            <ListTableRefresher clusters={rows.map(s => s.name)}/>
                            <ClusterDeploy keeper={keeper} database={KeeperPluginOptions[keeper].dbPlugin}/>
                            <ClusterDetect keeper={keeper} database={KeeperPluginOptions[keeper].dbPlugin}/>
                            <ListClusterAdd onClick={() => setShowNewElement(true)} disabled={showNewElement}/>
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
        return (
            <TableRow>
                <TableCell colSpan={3}>
                    {search ? (
                        <AlertCentered text={"There are no clusters that match your filter"}/>
                    ) : (
                        <ListEmptyInfo
                            onAddManually={() => setShowNewElement(true)}
                            disabledAddManually={showNewElement}
                        />
                    )}
                </TableCell>
            </TableRow>
        )
    }
}
