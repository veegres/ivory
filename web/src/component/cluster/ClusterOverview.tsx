import {Button, Grid, Radio, Table, TableCell, TableHead, TableRow} from "@mui/material";
import {useMutation, useQuery, useQueryClient} from "react-query";
import {nodeApi} from "../../app/api";
import {ErrorAlert} from "../view/ErrorAlert";
import {TableBody} from "../view/TableBody";
import React, {useState} from "react";
import {TableCellLoader} from "../view/TableCellLoader";
import {AxiosError} from "axios";
import {NodeColor} from "../../app/utils";
import {Node} from "../../app/types";
import {AlertDialog} from "../view/AlertDialog";
import {useStore} from "../../provider/StoreProvider";

const SX = {
    table: {'tr:last-child td': {border: 0}},
    clickable: {cursor: "pointer"},
    actionCell: {width: "50px"},
}

type AlertDialogState = {open: boolean, title: string, content: string, onAgree: () => void}

export function ClusterOverview() {
    const initAlertDialog = {open: false, title: '', content: '', onAgree: () => {}}
    const [alertDialog, setAlertDialog] = useState<AlertDialogState>(initAlertDialog)
    const { setStore, store: { activeCluster: { name: cluster }, activeNode } } = useStore()

    const clusterState = useQuery<Node[], AxiosError>(
        ["node/cluster", cluster],
        {retry: 0, refetchOnMount: false}
    )
    const queryClient = useQueryClient();
    const switchoverNode = useMutation(nodeApi.switchover, {
        onSuccess: async () => await queryClient.refetchQueries(['node/cluster', cluster])
    })
    const reinitNode = useMutation(nodeApi.reinitialize, {
        onSuccess: async () => await queryClient.refetchQueries(['node/cluster', cluster])
    })

    const {data: members, isLoading, isFetching, error} = clusterState
    if (error) return <ErrorAlert error={error}/>

    return (
        <>
            <AlertDialog {...alertDialog} onClose={() => setAlertDialog({...alertDialog, open: false})}/>
            <Table size="small" sx={SX.table}>
                <TableHead>
                    <TableRow>
                        <TableCell sx={SX.actionCell}/>
                        <TableCell>Role</TableCell>
                        <TableCell>Node</TableCell>
                        <TableCell>Host</TableCell>
                        <TableCell>State</TableCell>
                        <TableCell>Lag</TableCell>
                        <TableCellLoader sx={SX.actionCell} isFetching={(isFetching || switchoverNode.isLoading || reinitNode.isLoading)}/>
                    </TableRow>
                </TableHead>
                <TableBody isLoading={isLoading} cellCount={7}>
                    {renderContent()}
                </TableBody>
            </Table>
        </>
    )

    function renderContent() {
        if (!members) return <ErrorAlert error={"No data"}/>

        return members.map((node) => {
            const {api_domain, name, port, host, role, state, lag, isLeader} = node
            const isChecked = activeNode === api_domain
            return (
                <TableRow sx={SX.clickable} key={host} onClick={handleCheck(isChecked, api_domain)}>
                    <TableCell><Radio checked={isChecked} size={"small"}/></TableCell>
                    <TableCell sx={{color: NodeColor[role]}}>{role.toUpperCase()}</TableCell>
                    <TableCell>{name}</TableCell>
                    <TableCell>{host}:{port}</TableCell>
                    <TableCell>{state}</TableCell>
                    <TableCell>{lag}</TableCell>
                    <TableCell align="right" onClick={(e) => e.stopPropagation()}>
                        <Grid container justifyContent="flex-end" alignItems="center">
                            <Grid item>
                                {!isLeader ? (
                                    <Button
                                        color="primary"
                                        disabled={reinitNode.isLoading || switchoverNode.isLoading}
                                        onClick={() => handleReinit(api_domain)}>
                                        Reinit
                                    </Button>
                                ) : (
                                    <Button
                                        color="secondary"
                                        disabled={switchoverNode.isLoading}
                                        onClick={() => handleSwitchover(api_domain, name)}>
                                        Switchover
                                    </Button>
                                )}
                            </Grid>
                        </Grid>
                    </TableCell>
                </TableRow>
            )
        })
    }

    function handleCheck(checked: boolean, domain: string) {
        return () => setStore({ activeNode: checked ? '' : domain })
    }

    function handleSwitchover(node: string, leader: string) {
        setAlertDialog({
            open: true,
            title: `Switchover [${node}]`,
            content: 'Are you sure that you want to do Switchover? It will change the leader of your cluster.',
            onAgree: () => switchoverNode.mutate({node, leader})
        })
    }

    function handleReinit(node: string) {
        setAlertDialog({
            open: true,
            title: `Reinitialization [${node}]`,
            content: 'Are you sure that you want to do Reinit? It will erase all node data and will download it from scratch.',
            onAgree: () => reinitNode.mutate(node)
        })
    }
}
