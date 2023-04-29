import {ErrorSmart} from "../../view/box/ErrorSmart";
import {Box, Table, TableCell, TableHead, TableRow} from "@mui/material";
import {TableCellLoader} from "../../view/table/TableCellLoader";
import {TableBody} from "../../view/table/TableBody";
import {ListRowNew} from "./ListRowNew";
import {ListRow} from "./ListRow";
import {useState} from "react";
import {InfoAlert} from "../../view/box/InfoAlert";
import {AddIconButton} from "../../view/button/IconButtons";
import {Cluster} from "../../../type/cluster";
import {SxPropsMap} from "../../../type/common";
import scroll from "../../../style/scroll.module.css"
import {useAppearance} from "../../../provider/AppearanceProvider";
import {ListCreateAuto} from "./ListCreateAuto";

const SX: SxPropsMap = {
    box: {overflowY: "scroll"},
    table: {"tr:last-child td": {border: 0}, bgcolor: "inherit"},
    nameCell: {width: "220px"},
    buttonCell: {width: "1%"},
}

type Props = {
    rows: [string, Cluster][],
    isLoading: boolean,
    isFetching: boolean,
    selected: boolean,
    error: any,
}

export function ListTable(props: Props) {
    const {state: {mode}} = useAppearance()
    const {rows, error, isFetching, isLoading, selected} = props
    const [showNewElement, setShowNewElement] = useState(false)
    const [editNode, setEditNode] = useState("")

    if (error) return <ErrorSmart error={error}/>
    const header = {light: "#fff", dark: "#282828"}
    const bgcolor = header[mode] ?? "inherit"

    return (
        <Box sx={SX.box} className={scroll.tiny} maxHeight={selected ? "25vh" : "60vh"}>
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
                    {renderRows()}
                    {renderEmpty()}
                </TableBody>
            </Table>
        </Box>
    )

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
        if (isLoading || showNewElement || rows.length) return

        return (
            <TableRow>
                <TableCell colSpan={3}>
                    <InfoAlert text={"Please, add a cluster to get started"}/>
                </TableCell>
            </TableRow>
        )
    }
}
