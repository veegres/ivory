import {IconButton, Table, TableCell, TableHead, TableRow, Tooltip} from "@mui/material";
import { useQuery } from "react-query";
import {clusterApi} from "../../app/api";
import {ClustersRow} from "./ClustersRow";
import React, {useState} from "react";
import {Error} from "../view/Error";
import {TableBodyLoading} from "../view/TableBodyLoading";
import {TableCellFetching} from "../view/TableCellFetching";
import {AxiosError} from "axios";
import {Add} from "@mui/icons-material";

export function Clusters() {
    const [editNode, setEditNode] = useState('')
    const [showNewElement, setShowNewElement] = useState(false)
    const { data: clusterMap, isLoading, isFetching, isError, error } = useQuery('cluster/list', clusterApi.list)
    if (isError) return <Error error={error as AxiosError} />

    const closeNewElement = () => setShowNewElement(false)

    return (
        <Table size="small" sx={{ 'tr:last-child td': { border: 0 } }}>
            <TableHead>
                <TableRow>
                    <TableCell>Cluster Name</TableCell>
                    <TableCell>Nodes</TableCell>
                    <TableCellFetching isFetching={isFetching && !isLoading}>
                        <Tooltip title={'Add new cluster'} disableInteractive>
                            <IconButton disabled={showNewElement} onClick={() => setShowNewElement(true)}>
                                <Add />
                            </IconButton>
                        </Tooltip>
                    </TableCellFetching>
                </TableRow>
            </TableHead>
            <TableBodyLoading isLoading={isLoading} cellCount={3}>
                {Object.entries(clusterMap ?? {}).map(([name, nodes]) => {
                    const isReadOnly = name !== editNode
                    const toggleEdit = () => setEditNode(isReadOnly ? name : '')
                    const edit = { isReadOnly, toggleEdit, closeNewElement }
                    return <ClustersRow key={name} nodes={nodes} name={name} edit={edit} />
                })}
                {showNewElement ? <ClustersRow edit={{ closeNewElement }} /> : null}
            </TableBodyLoading>
        </Table>
    )
}
