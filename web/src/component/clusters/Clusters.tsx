import {IconButton, Table, TableCell, TableHead, TableRow, Tooltip} from "@mui/material";
import { useQuery } from "react-query";
import {clusterApi} from "../../app/api";
import {ClustersRow} from "./ClustersRow";
import React, {useState} from "react";
import {Error} from "../view/Error";
import {TableBodySkeleton} from "../view/TableBodySkeleton";
import {TableHeaderLoader} from "../view/TableHeaderLoader";
import {AxiosError} from "axios";
import {Add} from "@mui/icons-material";

const SX = {
    table: { 'tr:last-child td': { border: 0 } }
}

export function Clusters() {
    const [editNode, setEditNode] = useState('')
    const [showNewElement, setShowNewElement] = useState(false)
    const { data: clusterMap, isLoading, isFetching, isError, error } = useQuery('cluster/list', clusterApi.list)
    if (isError) return <Error error={error as AxiosError} />

    const closeNewElement = () => setShowNewElement(false)

    return (
        <Table size="small" sx={SX.table}>
            <TableHead>
                <TableRow>
                    <TableCell>Cluster Name</TableCell>
                    <TableCell>Nodes</TableCell>
                    <TableHeaderLoader isFetching={isFetching && !isLoading}>
                        <Tooltip title={'Add new cluster'} disableInteractive>
                            <IconButton disabled={showNewElement} onClick={() => setShowNewElement(true)}>
                                <Add />
                            </IconButton>
                        </Tooltip>
                    </TableHeaderLoader>
                </TableRow>
            </TableHead>
            <TableBodySkeleton isLoading={isLoading} cellCount={3}>
                {Object.entries(clusterMap ?? {}).map(([name, nodes]) => {
                    const isReadOnly = name !== editNode
                    const toggleEdit = () => setEditNode(isReadOnly ? name : '')
                    const edit = { isReadOnly, toggleEdit, closeNewElement }

                    return <ClustersRow key={name} nodes={nodes} name={name} edit={edit} />
                })}
                {showNewElement ? <ClustersRow edit={{ closeNewElement }} /> : null}
            </TableBodySkeleton>
        </Table>
    )
}
