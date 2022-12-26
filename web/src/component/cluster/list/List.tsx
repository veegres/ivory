import {IconButton, Table, TableCell, TableHead, TableRow, Tooltip} from "@mui/material";
import {useQuery} from "@tanstack/react-query";
import {clusterApi} from "../../../app/api";
import {ListRow} from "./ListRow";
import React, {useState} from "react";
import {ErrorAlert} from "../../view/ErrorAlert";
import {TableBody} from "../../view/TableBody";
import {TableCellLoader} from "../../view/TableCellLoader";
import {Add} from "@mui/icons-material";
import {Block} from "../../view/Block";
import {ListRowNew} from "./ListRowNew";

const SX = {
    table: {'tr:last-child td': {border: 0}},
    nameCell: {width: '220px'},
    buttonCell: {width: '1%'}
}

export function List() {
    const [editNode, setEditNode] = useState('')
    const [showNewElement, setShowNewElement] = useState(false)
    const {data: clusterMap, isLoading, isFetching, isError, error} = useQuery(["cluster/list"], clusterApi.list)

    return (
        <Block>
            {renderContent()}
        </Block>
    )

    function renderContent() {
        if (isError) return <ErrorAlert error={error}/>

        return (
            <Table size={"small"} sx={SX.table}>
                <TableHead>
                    <TableRow>
                        <TableCell sx={SX.nameCell}>Cluster Name</TableCell>
                        <TableCell>Instances</TableCell>
                        <TableCellLoader sx={SX.buttonCell} isFetching={isFetching && !isLoading}>
                            <Tooltip title={'Add new cluster'} disableInteractive>
                                <IconButton disabled={showNewElement} onClick={() => setShowNewElement(true)}>
                                    <Add/>
                                </IconButton>
                            </Tooltip>
                        </TableCellLoader>
                    </TableRow>
                </TableHead>
                <TableBody isLoading={isLoading} cellCount={3}>
                    {renderRows()}
                    <ListRowNew show={showNewElement} close={() => setShowNewElement(false)} />
                </TableBody>
            </Table>
        )
    }

    function renderRows() {
        return Object.entries(clusterMap ?? {}).map(([name, cluster]) => {
            const editable = name === editNode
            const toggle = () => setEditNode(editable ? '' : name)

            return <ListRow key={name} cluster={cluster} name={name} editable={editable} toggle={toggle}/>
        })
    }
}
