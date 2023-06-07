import {ErrorSmart} from "../../view/box/ErrorSmart";
import {Box, Table, TableCell, TableHead, TableRow} from "@mui/material";
import {TableCellLoader} from "../../view/table/TableCellLoader";
import {TableBody} from "../../view/table/TableBody";
import {ListRowNew} from "./ListRowNew";
import {ListRow} from "./ListRow";
import {useMemo, useState} from "react";
import {InfoAlert} from "../../view/box/InfoAlert";
import {AddIconButton} from "../../view/button/IconButtons";
import {ClusterMap} from "../../../type/cluster";
import {SxPropsMap} from "../../../type/common";
import scroll from "../../../style/scroll.module.css"
import {useAppearance} from "../../../provider/AppearanceProvider";
import {ListCreateAuto} from "./ListCreateAuto";
import {useStore} from "../../../provider/StoreProvider";

const SX: SxPropsMap = {
    box: {overflowY: "scroll"},
    table: {"tr:last-child td": {border: 0}, bgcolor: "inherit"},
    nameCell: {width: "220px"},
    buttonCell: {width: "1%"},
}

type Props = {
    map?: ClusterMap,
    isLoading: boolean,
    isFetching: boolean,
    error: any,
}

export function ListTable(props: Props) {
    const {store: {activeCluster}} = useStore()
    const {state: {mode}} = useAppearance()
    const {map, error, isFetching, isLoading} = props
    const [showNewElement, setShowNewElement] = useState(false)
    const [editNode, setEditNode] = useState("")

    const rows = useMemo(() => Object.entries(map ?? {}), [map])

    if (error) return <ErrorSmart error={error}/>
    const header = {light: "#fff", dark: "#282828"}
    const bgcolor = header[mode] ?? "inherit"

    return (
        <Box sx={SX.box} className={scroll.tiny} maxHeight={!!activeCluster ? "25vh" : "60vh"}>
            <Table size={"small"} sx={SX.table} stickyHeader>
                <TableHead>
                    <TableRow>
                        <TableCell sx={{bgcolor, ...SX.nameCell}}>Cluster Name</TableCell>
                        <TableCell sx={{bgcolor}}>Instances</TableCell>
                        <TableCellLoader sx={{bgcolor, ...SX.buttonCell}} isFetching={isFetching && !isLoading}>
                            <ListCreateAuto />
                            <AddIconButton onClick={() => setShowNewElement(true)} disabled={showNewElement}/>
                        </TableCellLoader>
                    </TableRow>
                </TableHead>
                <TableBody isLoading={isLoading} cellCount={3}>
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
                <ListRow
                    key={name}
                    cluster={cluster}
                    editable={editable}
                    toggle={toggle}
                />
            )
        })
    }

    function renderEmpty() {
        if (isLoading || showNewElement || rows.length || activeCluster) return

        return (
            <TableRow>
                <TableCell colSpan={3}>
                    <InfoAlert text={"There is no clusters yet. You can add them manually or by auto detection."}/>
                </TableCell>
            </TableRow>
        )
    }
}
