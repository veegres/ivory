import {Table, TableCell, TableHead, TableRow} from "@mui/material";
import { useQuery } from "react-query";
import {clusterApi} from "../../app/api";
import {ClusterListRow} from "./ClusterListRow";
import React, {Dispatch} from "react";
import {Error} from "../view/Error";
import {TableBodyLoading} from "../view/TableBodyLoading";
import {TableCellFetching} from "../view/TableCellFetching";
import {AxiosError} from "axios";

export function ClusterList({ setNode }: { setNode: Dispatch<string> }) {
    const { data: clusterMap, isLoading, isFetching, isError, error } = useQuery('cluster/list', clusterApi.list)
    if (isError) return <Error error={error as AxiosError} />

    return (
        <Table size="small" sx={{ 'tr:last-child td': { border: 0 } }}>
            <TableHead>
                <TableRow>
                    <TableCell>Name</TableCell>
                    <TableCell>Nodes</TableCell>
                    <TableCellFetching isFetching={isFetching && !isLoading} />
                </TableRow>
            </TableHead>
            <TableBodyLoading isLoading={isLoading} cellCount={3}>
                <Content />
            </TableBodyLoading>
        </Table>
    )

    function Content() {
        if (!clusterMap) return null

        return (
            <>
                {Object.entries(clusterMap).map(([name, nodes]) => (
                    <ClusterListRow key={name} nodes={nodes} name={name} setNode={setNode} />
                ))}
                <ClusterListRow />
            </>
        )
    }
}
