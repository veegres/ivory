import {Table, TableBody, TableCell, TableHead, TableRow} from "@mui/material";
import { useQuery } from "react-query";
import {clusterApi} from "../../app/api";
import {ClusterCreateComponent} from "./ClusterCreateComponent";
import React from "react";

export function ClusterListComponent() {
    const { data: clusterList } = useQuery('cluster/list', clusterApi.list)
    if (!clusterList || clusterList.length === 0) return null;
    return (
        <Table size="small" sx={{ 'tr:last-child td': { border: 0 } }}>
            <TableHead>
                <TableRow>
                    <TableCell sx={{ width: '15%' }}>Name</TableCell>
                    <TableCell>Nodes</TableCell>
                    <TableCell sx={{ width: '115px' }} />
                </TableRow>
            </TableHead>
            <TableBody>
                {clusterList.map(cluster => (
                    <ClusterCreateComponent key={cluster.name} nodes={cluster.nodes} name={cluster.name} />
                ))}
                <ClusterCreateComponent nodes={[]} />
            </TableBody>
        </Table>
    )
}
