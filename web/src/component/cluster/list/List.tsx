import {Box, IconButton, Table, TableCell, TableHead, TableRow, Tooltip} from "@mui/material";
import {useQuery} from "@tanstack/react-query";
import {clusterApi} from "../../../app/api";
import {ListRow} from "./ListRow";
import React, {useState} from "react";
import {ErrorAlert} from "../../view/ErrorAlert";
import {TableBody} from "../../view/TableBody";
import {TableCellLoader} from "../../view/TableCellLoader";
import {Add} from "@mui/icons-material";
import {PageBlock} from "../../view/PageBlock";
import {ListRowNew} from "./ListRowNew";
import {InfoAlert} from "../../view/InfoAlert";

const SX = {
    table: {"tr:last-child td": {border: 0}},
    nameCell: {width: "220px"},
    buttonCell: {width: "1%"},
    padding: {padding: "10px 20px"},
}

export function List() {
    const [editNode, setEditNode] = useState("")
    const [showNewElement, setShowNewElement] = useState(false)
    const query = useQuery(["cluster/list"], clusterApi.list)
    const {data: clusterMap, isLoading, isFetching, isError, error} = query
    const rows = Object.entries(clusterMap ?? {})

    return (
        <PageBlock>
            {renderContent()}
        </PageBlock>
    )

    function renderContent() {
        if (isError) return <ErrorAlert error={error}/>

        return (
            <>
                <Table size={"small"} sx={SX.table}>
                    <TableHead>
                        <TableRow>
                            <TableCell sx={SX.nameCell}>Cluster Name</TableCell>
                            <TableCell>Instances</TableCell>
                            <TableCellLoader sx={SX.buttonCell} isFetching={isFetching && !isLoading}>
                                <Tooltip title={"Add new cluster"} disableInteractive>
                                    <IconButton disabled={showNewElement} onClick={() => setShowNewElement(true)}>
                                        <Add/>
                                    </IconButton>
                                </Tooltip>
                            </TableCellLoader>
                        </TableRow>
                    </TableHead>
                    <TableBody isLoading={isLoading} cellCount={3}>
                        {renderRows()}
                        <ListRowNew show={showNewElement} close={() => setShowNewElement(false)}/>
                    </TableBody>

                </Table>
                {!showNewElement && !rows.length && (
                    <Box sx={SX.padding}>
                        <InfoAlert text={"Please, add a cluster to get started"}/>
                    </Box>
                )}
            </>
        )
    }

    function renderRows() {
        return rows.map(([name, cluster]) => {
            const editable = name === editNode
            const toggle = () => setEditNode(editable ? "" : name)

            return <ListRow key={name} cluster={cluster} name={name} editable={editable} toggle={toggle}/>
        })
    }
}
