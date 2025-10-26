import {Box, Table, TableCell, TableHead, TableRow} from "@mui/material"
import {useMemo, useState} from "react"

import {ClusterMap} from "../../../../api/cluster/type"
import {SxPropsMap} from "../../../../app/type"
import {SxPropsFormatter} from "../../../../app/utils"
import {useStore} from "../../../../provider/StoreProvider"
import scroll from "../../../../style/scroll.module.css"
import {ErrorSmart} from "../../../view/box/ErrorSmart"
import {InfoAlert} from "../../../view/box/InfoAlert"
import {AddIconButton} from "../../../view/button/IconButtons"
import {TableBody} from "../../../view/table/TableBody"
import {TableCellLoader} from "../../../view/table/TableCellLoader"
import {ListCreateAuto} from "./ListCreateAuto"
import {ListRow} from "./ListRow"
import {ListRowNew} from "./ListRowNew"

const SX: SxPropsMap = {
    box: {overflowY: "scroll"},
    table: {"tr:last-child td": {border: 0}, "tr td, th": {padding: "5px 10px"}},
    nameCell: {width: "220px"},
    buttonCell: {width: "130px"},
}

type Props = {
    map?: ClusterMap,
    pending: boolean,
    fetching: boolean,
    error: any,
}

export function ListTable(props: Props) {
    const activeCluster = useStore(s => s.activeCluster)
    const {map, error, fetching, pending} = props
    const [showNewElement, setShowNewElement] = useState(false)
    const [editNode, setEditNode] = useState("")

    const rows = useMemo(() => Object.entries(map ?? {}), [map])

    if (error) return <ErrorSmart error={error}/>

    return (
        <Box sx={SX.box} className={scroll.tiny} maxHeight={activeCluster ? "25vh" : "60vh"}>
            <Table size={"small"} sx={SX.table} stickyHeader>
                <TableHead>
                    <TableRow>
                        <TableCell sx={[SxPropsFormatter.style.paper, SX.nameCell]}>Cluster Name</TableCell>
                        <TableCell sx={SxPropsFormatter.style.paper}>Instances</TableCell>
                        <TableCellLoader sx={[SxPropsFormatter.style.paper, SX.buttonCell]} isFetching={fetching && !pending}>
                            <ListCreateAuto />
                            <AddIconButton tooltip={"Add Cluster Manually"} onClick={() => setShowNewElement(true)} disabled={showNewElement}/>
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
        if (!activeCluster || !map) return
        if (map[activeCluster.cluster.name] !== undefined) return
        return (
            <ListRow cluster={activeCluster.cluster} editable={false}/>
        )
    }

    function renderRows() {
        return rows.map(([name, cluster]) => {
            const editable = name === editNode
            const toggle = () => setEditNode(editable ? "" : name)

            return (
                <ListRow key={name} cluster={cluster} editable={editable} toggle={toggle}/>
            )
        })
    }

    function renderEmpty() {
        if (pending || showNewElement || rows.length || activeCluster) return

        return (
            <TableRow>
                <TableCell colSpan={3}>
                    <InfoAlert text={"There is no clusters yet. You can add them manually or by auto detection."}/>
                </TableCell>
            </TableRow>
        )
    }
}
