import {Table, TableCell, TableHead, TableRow} from "@mui/material";
import { useQuery } from "react-query";
import {clusterApi} from "../../app/api";
import {ClusterListRow} from "./ClusterListRow";
import React, {Dispatch} from "react";
import {Error} from "../view/Error";
import {TableBodyLoading} from "../view/TableBodyLoading";
import {TableCellFetching} from "../view/TableCellFetching";

export function ClusterList({ setNode }: { setNode: Dispatch<string> }) {
    const { data: clusterList, isLoading, isFetching, isError, error } = useQuery('cluster/list', clusterApi.list)
    if (isError) return <Error error={error} />

    return (
        <Table size="small" sx={{ 'tr:last-child td': { border: 0 } }}>
            <TableHead>
                <TableRow>
                    <TableCell>Name</TableCell>
                    <TableCell>Nodes</TableCell>
                    <TableCellFetching isFetching={isFetching} />
                </TableRow>
            </TableHead>
            <TableBodyLoading isLoading={isLoading} cellCount={3}>
                <Content />
            </TableBodyLoading>
        </Table>
    )

    function Content() {
        if (!clusterList || clusterList.length === 0) return null;

        return (
            <>
                {clusterList.map(cluster => (
                    <ClusterListRow key={cluster.name} nodes={cluster.nodes} name={cluster.name} setNode={setNode} />
                ))}
                <ClusterListRow />
            </>
        )
    }
}
