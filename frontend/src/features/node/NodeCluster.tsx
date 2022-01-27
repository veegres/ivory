import {Button, Grid, Table, TableCell, TableHead, TableRow} from "@mui/material";
import {OpenInNew} from "@mui/icons-material";
import { useQuery } from "react-query";
import {nodeApi} from "../../app/api";
import {Error} from "../view/Error";
import {TableBodyLoading} from "../view/TableBodyLoading";
import React from "react";
import {TableCellFetching} from "../view/TableCellFetching";
import {AxiosError} from "axios";

const SX = {
    tableLastChildRow: { 'tr:last-child td': { border: 0 } }
}

export function NodeCluster({ node }: { node: string }) {
    const { data: members, isLoading, isFetching, isError, error } = useQuery(['node/cluster', node], () => nodeApi.cluster(node))
    if (isError) return <Error error={error as AxiosError} />

    return (
        <Table size="small" sx={SX.tableLastChildRow}>
            <TableHead>
                <TableRow>
                    <TableCell>Node</TableCell>
                    <TableCell>Role</TableCell>
                    <TableCell>Lag</TableCell>
                    <TableCellFetching isFetching={isFetching && !isLoading} />
                </TableRow>
            </TableHead>
            <TableBodyLoading isLoading={isLoading} cellCount={4}>
                <Content />
            </TableBodyLoading>
        </Table>
    )

    function Content() {
        if (!members || members.length === 0) return null

        return (
            <>
                {members.map((node) => (
                    <TableRow key={node.host}>
                        <TableCell>{node.name}</TableCell>
                        <TableCell>{node.role.toUpperCase()}</TableCell>
                        <TableCell>{node.lag}</TableCell>
                        <TableCell align="right">
                            <Grid container justifyContent="flex-end" alignItems="center">
                                <Grid item>{node.role === "leader" ? <Button color="secondary">Switchover</Button> : null}</Grid>
                                <Grid item><Button color="primary">Reinit</Button></Grid>
                                <Grid item><Button href={node.api_url} target="_blank"><OpenInNew /></Button></Grid>
                            </Grid>
                        </TableCell>
                    </TableRow>
                ))}
            </>
        )
    }
}
