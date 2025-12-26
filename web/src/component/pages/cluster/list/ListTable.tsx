import {Box, Table, TableCell, TableHead, TableRow} from "@mui/material"
import {useMemo, useState} from "react"

import {ClusterApi} from "../../../../api/cluster/router"
import {Cluster} from "../../../../api/cluster/type"
import {Permission} from "../../../../api/permission/type"
import {SxPropsMap} from "../../../../app/type"
import {SxPropsFormatter} from "../../../../app/utils"
import {useStore} from "../../../../provider/StoreProvider"
import scroll from "../../../../style/scroll.module.css"
import {AlertCentered} from "../../../view/box/AlertCentered"
import {AddIconButton} from "../../../view/button/IconButtons"
import {TableBody} from "../../../view/table/TableBody"
import {TableCellLoader} from "../../../view/table/TableCellLoader"
import {Access} from "../../../widgets/access/Access"
import {Refresher} from "../../../widgets/refresher/Refresher"
import {ListCreateAuto} from "./ListCreateAuto"
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
                        <TableCell sx={SxPropsFormatter.style.paper}>Cluster Name</TableCell>
                        <TableCellLoader
                            sx={SxPropsFormatter.style.paper}
                            label={"Instances"}
                            colSpan={2}
                            loading={fetching && !pending}
                        >
                            <Box sx={SX.refresh}>
                                <Refresher queryKeys={[ClusterApi.list.key(), ClusterApi.overview.key()]}/>
                            </Box>
                            <ListCreateAuto/>
                            <Access permission={Permission.ManageClusterUpdate}>
                                <AddIconButton
                                    tooltip={"Add Cluster Manually"}
                                    onClick={() => setShowNewElement(true)}
                                    disabled={showNewElement}
                                />
                            </Access>
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
        if (!activeCluster || !list.length) return
        if (list.some(e => e.name === activeCluster.cluster.name)) return
        return (
            <ListRow cluster={activeCluster.cluster} editable={false}/>
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
                    <AlertCentered text={"There are no clusters yet. You can add them manually or by auto detection."}/>
                </TableCell>
            </TableRow>
        )
    }
}
