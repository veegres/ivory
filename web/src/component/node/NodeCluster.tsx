import {Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, Grid, Table, TableCell, TableHead, TableRow} from "@mui/material";
import {useMutation, useQuery, useQueryClient} from "react-query";
import {nodeApi} from "../../app/api";
import {Error} from "../view/Error";
import {TableBodySkeleton} from "../view/TableBodySkeleton";
import React, {useState} from "react";
import {TableHeaderLoader} from "../view/TableHeaderLoader";
import {AxiosError} from "axios";
import {nodeColor} from "../../app/utils";

const SX = {
    tableLastChildRow: { 'tr:last-child td': { border: 0 } }
}

type Props = { node: string }
type AlertDialog = { isOpen: boolean, title?: string, content?: string, agree?: () => void }

export function NodeCluster({ node }: Props) {
    const [alertDialog, setAlertDialog] = useState<AlertDialog>({ isOpen: false })

    const { data: members, isLoading, isFetching, isError, error } = useQuery(
        ['node/cluster', node],
        () => nodeApi.cluster(node),
        { refetchInterval: 5000 }
    )

    const queryClient = useQueryClient();
    const switchoverNode = useMutation(nodeApi.switchover, {
        onSuccess: async () => await queryClient.refetchQueries(['node/cluster', node])
    })
    const reinitNode = useMutation(nodeApi.reinitialize, {
        onSuccess: async () => await queryClient.refetchQueries(['node/cluster', node])
    })

    if (isError) return <Error error={error as AxiosError} />

    return (
        <>
            {renderAlertDialog()}
            <Table size="small" sx={SX.tableLastChildRow}>
                <TableHead>
                    <TableRow>
                        <TableCell>Node</TableCell>
                        <TableCell>Role</TableCell>
                        <TableCell>State</TableCell>
                        <TableCell>Lag</TableCell>
                        <TableHeaderLoader isFetching={(isFetching || switchoverNode.isLoading || reinitNode.isLoading) && !isLoading} />
                    </TableRow>
                </TableHead>
                <TableBodySkeleton isLoading={isLoading} cellCount={5}>
                    {renderContent()}
                </TableBodySkeleton>
            </Table>
        </>
    )

    function renderContent() {
        if (!members || members.length === 0) return null

        return (
            <>
                {members.map((node) => {
                    const nodePublicHost = node.api_url.split('/')[2]
                    const isLeader = node.role === "leader"
                    return (
                        <TableRow key={node.host}>
                            <TableCell>{node.name}</TableCell>
                            <TableCell sx={{ color: nodeColor[node.role] }}>{node.role.toUpperCase()}</TableCell>
                            <TableCell>{node.state}</TableCell>
                            <TableCell>{node.lag}</TableCell>
                            <TableCell align="right">
                                <Grid container justifyContent="flex-end" alignItems="center">
                                    <Grid item>
                                        {!isLeader ? null : (
                                            <Button
                                                color="secondary"
                                                disabled={switchoverNode.isLoading}
                                                onClick={() => handleSwitchover(nodePublicHost, node.name)}>
                                                Switchover
                                            </Button>
                                        )}
                                    </Grid>
                                    <Grid item>
                                        {isLeader ? null : (
                                            <Button
                                                color="primary"
                                                disabled={reinitNode.isLoading || switchoverNode.isLoading}
                                                onClick={() => handleReinit(nodePublicHost)}>
                                                Reinit
                                            </Button>
                                        )}
                                    </Grid>
                                </Grid>
                            </TableCell>
                        </TableRow>
                    )
                })}
            </>
        )
    }

    function renderAlertDialog() {
        const { isOpen, content, title, agree } = alertDialog
        if (!content || !title || !agree) return null
        return (
            <Dialog
                open={isOpen}
                onClose={handleClose}
                aria-labelledby="alert-dialog-title"
                aria-describedby="alert-dialog-description"
            >
                <DialogTitle id="alert-dialog-title">
                    {title}
                </DialogTitle>
                <DialogContent>
                    <DialogContentText id="alert-dialog-description">
                        {content}
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleClose}>No, please, No...</Button>
                    <Button onClick={() => { agree(); handleClose() }} autoFocus>Yes, move on!</Button>
                </DialogActions>
            </Dialog>
        )
    }

    function handleClose() {
        return setAlertDialog({ isOpen: false, content: undefined, title: undefined })
    }

    function handleSwitchover(node: string, leader: string) {
        setAlertDialog({
            isOpen: true,
            title: 'Switchover',
            content: 'Are you sure that you want to do Switchover? It will change the leader of your cluster.',
            agree: () => switchoverNode.mutate({ node, leader })
        })
    }

    function handleReinit(node: string) {
        setAlertDialog({
            isOpen: true,
            title: 'Reinitialization',
            content: 'Are you sure that you want to do Reinit? It will erase all node data and will download it from scratch.',
            agree: () => reinitNode.mutate(node)
        })
    }
}
