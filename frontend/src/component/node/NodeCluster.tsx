import {Button, Grid, Table, TableCell, TableHead, TableRow} from "@mui/material";
import {OpenInNew} from "@mui/icons-material";
import {useMutation, useQuery, useQueryClient} from "react-query";
import {nodeApi} from "../../app/api";
import {Error} from "../view/Error";
import {TableBodySkeleton} from "../view/TableBodySkeleton";
import React from "react";
import {TableCellFetching} from "../view/TableCellFetching";
import {AxiosError} from "axios";

const SX = {
    tableLastChildRow: { 'tr:last-child td': { border: 0 } }
}

export function NodeCluster({ node }: { node: string }) {
    const { data: members, isLoading, isFetching, isError, error } = useQuery(['node/cluster', node], () => nodeApi.cluster(node))

    const queryClient = useQueryClient();
    const switchoverNode = useMutation(nodeApi.switchover, {
        onSuccess: async () => {
            await queryClient.refetchQueries(['node/cluster', node])
        }
    })
    const reinitNode = useMutation(nodeApi.reinitialize, {
        onSuccess: async () => {
            await queryClient.refetchQueries(['node/cluster', node])
        }
    })

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
            <TableBodySkeleton isLoading={isLoading} cellCount={4}>
                <Content />
            </TableBodySkeleton>
        </Table>
    )

    function Content() {
        if (!members || members.length === 0) return null

        return (
            <>
                {members.map((node) => {
                    const nodePublicHost = node.api_url.split('/')[2]
                    const isLeader = node.role === "leader"
                    return (
                        <TableRow key={node.host}>
                            <TableCell>{node.name}</TableCell>
                            <TableCell>{node.role.toUpperCase()}</TableCell>
                            <TableCell>{node.lag}</TableCell>
                            <TableCell align="right">
                                <Grid container justifyContent="flex-end" alignItems="center">
                                    <Grid item>
                                        {!isLeader ? null : (
                                            <Button color="secondary" onClick={() => switchoverNode.mutate({ node: nodePublicHost, leader: node.name })}>
                                                Switchover
                                            </Button>
                                        )}
                                    </Grid>
                                    <Grid item>
                                        {isLeader ? null : (
                                            <Button color="primary" onClick={() => reinitNode.mutate(nodePublicHost)}>
                                                Reinit
                                            </Button>
                                        )}
                                    </Grid>
                                    <Grid item>
                                        <Button href={node.api_url} target="_blank"><OpenInNew /></Button>
                                    </Grid>
                                </Grid>
                            </TableCell>
                        </TableRow>
                    )
                })}
            </>
        )
    }
}
