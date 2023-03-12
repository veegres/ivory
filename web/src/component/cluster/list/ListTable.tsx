import {ErrorAlert} from "../../view/ErrorAlert";
import {Table, TableCell, TableHead, TableRow} from "@mui/material";
import {TableCellLoader} from "../../view/TableCellLoader";
import {TableBody} from "../../view/TableBody";
import {ListRowNew} from "./ListRowNew";
import {ListRow} from "./ListRow";
import {useState} from "react";
import {InfoAlert} from "../../view/InfoAlert";
import {AddIconButton, AutoIconButton} from "../../view/IconButtons";
import {Cluster} from "../../../type/cluster";
import {SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    table: {"tr:last-child td": {border: 0}},
    nameCell: {width: "220px"},
    buttonCell: {width: "1%"},
}

type Props = {
    rows: [string, Cluster][],
    isLoading: boolean,
    isFetching: boolean,
    error: any,
}

export function ListTable(props: Props) {
    const {rows, error, isFetching, isLoading} = props
    const [showNewElement, setShowNewElement] = useState(false)
    const [editNode, setEditNode] = useState("")

    if (error) return <ErrorAlert error={error}/>

    return (
        <Table size={"small"} sx={SX.table}>
            <TableHead>
                <TableRow>
                    <TableCell sx={SX.nameCell}>Cluster Name</TableCell>
                    <TableCell>Instances</TableCell>
                    <TableCellLoader sx={SX.buttonCell} isFetching={isFetching && !isLoading}>
                        <AutoIconButton onClick={() => {}} disabled/>
                        <AddIconButton onClick={() => setShowNewElement(true)} disabled={showNewElement}/>
                    </TableCellLoader>
                </TableRow>
            </TableHead>
            <TableBody isLoading={isLoading} cellCount={3}>
                {renderRows()}
                <ListRowNew show={showNewElement} close={() => setShowNewElement(false)}/>
                {renderEmpty()}
            </TableBody>
        </Table>
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
