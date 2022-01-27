import {Table, TableBody, TableCell, TableHead, TableRow} from "@mui/material";
import { useQuery } from "react-query";
import {clusterApi} from "../../app/api";
import {ClusterListRow} from "./ClusterListRow";
import React from "react";

export function ClusterList() {
    const { data: clusterList } = useQuery('cluster/list', clusterApi.list)
    if (!clusterList || clusterList.length === 0) return null;
    return (
        <Table size="small" sx={{ 'tr:last-child td': { border: 0 } }}>
            <TableHead>
                <TableRow>
                    <TableCell>Name</TableCell>
                    <TableCell>Nodes</TableCell>
                    <TableCell />
                </TableRow>
            </TableHead>
            <TableBody>
                {clusterList.map(cluster => (
                    <ClusterListRow key={cluster.name} nodes={cluster.nodes} name={cluster.name} />
                ))}
                <ClusterListRow />
            </TableBody>
        </Table>
    )
}
