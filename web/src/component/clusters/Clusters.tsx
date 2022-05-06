import {IconButton, Table, TableCell, TableHead, TableRow, Tooltip} from "@mui/material";
import {useQuery} from "react-query";
import {clusterApi} from "../../app/api";
import {ClustersRow} from "./ClustersRow";
import React, {useState} from "react";
import {Error} from "../view/Error";
import {TableBody} from "../view/TableBody";
import {TableCellLoader} from "../view/TableCellLoader";
import {AxiosError} from "axios";
import {Add} from "@mui/icons-material";
import {Block} from "../view/Block";

const SX = {
    table: {'tr:last-child td': {border: 0}},
    nameCell: {width: '220px'},
    buttonCell: {width: '1%'}
}

export function Clusters() {
    const [editNode, setEditNode] = useState('')
    const [showNewElement, setShowNewElement] = useState(false)
    const {data: clusterMap, isLoading, isFetching, isError, error} = useQuery('cluster/list', clusterApi.list)

    return (
        <Block>
            {renderContent()}
        </Block>
    )

    function renderContent() {
        if (isError) return <Error error={error as AxiosError}/>

        return (
            <Table size="small" sx={SX.table}>
                <TableHead>
                    <TableRow>
                        <TableCell sx={SX.nameCell}>Cluster Name</TableCell>
                        <TableCell>Nodes</TableCell>
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
                    {Object.entries(clusterMap ?? {}).map(([name, nodes]) => {
                        const isReadOnly = name !== editNode
                        const toggleEdit = () => setEditNode(isReadOnly ? name : '')
                        const edit = {isReadOnly, toggleEdit, closeNewElement}

                        return <ClustersRow key={name} nodes={nodes} name={name} edit={edit}/>
                    })}
                    {showNewElement ? <ClustersRow nodes={[]} name={''} edit={{closeNewElement}}/> : null}
                </TableBody>
            </Table>
        )
    }

    function closeNewElement() {
        setShowNewElement(false)
    }
}
