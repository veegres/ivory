import {Box, Button, Radio, Table, TableCell, TableHead, TableRow} from "@mui/material";
import {useMutation, useQuery} from "react-query";
import {nodeApi} from "../../app/api";
import {TableBody} from "../view/TableBody";
import React, {useState} from "react";
import {TableCellLoader} from "../view/TableCellLoader";
import {AxiosError} from "axios";
import {InstanceColor} from "../../app/utils";
import {InstanceMap} from "../../app/types";
import {AlertDialog} from "../view/AlertDialog";
import {useStore} from "../../provider/StoreProvider";
import {TapProps} from "./Cluster";

const SX = {
    table: {'tr:last-child td': {border: 0}},
    clickable: {cursor: "pointer"},
    actionCell: {width: "70px"},
    roleCell: {width: "110px"},
    buttonCell: {width: "160px"},
}

type AlertDialogState = {open: boolean, title: string, content: string, onAgree: () => void}
const initAlertDialog = {open: false, title: '', content: '', onAgree: () => {}}

export function ClusterOverview(props: TapProps) {
    const { cluster, instance } = props
    const [alertDialog, setAlertDialog] = useState<AlertDialogState>(initAlertDialog)
    const { setStore, store: { activeNode } } = useStore()

    const instanceMap = useQuery<InstanceMap, AxiosError>(
        ["node/cluster", instance.api_domain],
        {retry: 0, refetchOnMount: false}
    )
    const switchoverNode = useMutation(nodeApi.switchover, {
        onSuccess: async () => await instanceMap.refetch()
    })
    const reinitNode = useMutation(nodeApi.reinitialize, {
        onSuccess: async () => await instanceMap.refetch()
    })

    return (
        <>
            <AlertDialog {...alertDialog} onClose={() => setAlertDialog({...alertDialog, open: false})}/>
            <Table size={"small"} sx={SX.table}>
                <TableHead>
                    <TableRow>
                        <TableCell sx={SX.actionCell}/>
                        <TableCell sx={SX.roleCell}>Role</TableCell>
                        <TableCell align={"center"}>Patroni</TableCell>
                        <TableCell align={"center"}>Postgres</TableCell>
                        <TableCell align={"center"}>State</TableCell>
                        <TableCell align={"center"}>Lag</TableCell>
                        <TableCellLoader sx={SX.buttonCell} isFetching={(instanceMap.isFetching || switchoverNode.isLoading || reinitNode.isLoading)}/>
                    </TableRow>
                </TableHead>
                <TableBody isLoading={instanceMap.isLoading} cellCount={7}>
                    {renderContent()}
                </TableBody>
            </Table>
        </>
    )

    function renderContent() {
        return cluster.nodes.map((name) => {
            const isChecked = activeNode === name
            const element = instanceMap.data ? instanceMap.data[name] : undefined
            const { role, port, host, state, lag, api_domain, leader } = element ??
            { host: "-", port: "-", role: "unknown", api_domain: "-", lag: "-", leader: false, state: "-" }

            return (
                <TableRow sx={SX.clickable} key={name} onClick={handleCheck(isChecked, name)}>
                    <TableCell><Radio checked={isChecked} size={"small"}/></TableCell>
                    <TableCell sx={{color: InstanceColor[role]}}>{role.toUpperCase()}</TableCell>
                    <TableCell align={"center"}>{name}</TableCell>
                    <TableCell align={"center"}>{host}:{port}</TableCell>
                    <TableCell align={"center"}>{state}</TableCell>
                    <TableCell align={"center"}>{lag}</TableCell>
                    <TableCell align={"right"} onClick={(e) => e.stopPropagation()}>
                        {element && renderButton(api_domain, host, leader)}
                    </TableCell>
                </TableRow>
            )
        })
    }

    function renderButton(instance: string, host: string, leader: boolean) {
        return (
            <Box display={"flex"} justifyContent={"flex-end"} alignItems={"center"}>
                {!leader ? (
                    <Button
                        color={"primary"}
                        disabled={reinitNode.isLoading || switchoverNode.isLoading}
                        onClick={() => handleReinit(instance)}>
                        Reinit
                    </Button>
                ) : (
                    <Button
                        color={"secondary"}
                        disabled={switchoverNode.isLoading}
                        onClick={() => handleSwitchover(instance, host)}>
                        Switchover
                    </Button>
                )}
            </Box>
        )
    }

    function handleCheck(checked: boolean, domain: string) {
        return () => setStore({ activeNode: checked ? '' : domain })
    }

    function handleSwitchover(instance: string, host: string) {
        setAlertDialog({
            open: true,
            title: `Switchover [${instance}]`,
            content: 'Are you sure that you want to do Switchover? It will change the leader of your cluster.',
            onAgree: () => switchoverNode.mutate({node: instance, leader: host})
        })
    }

    function handleReinit(instance: string) {
        setAlertDialog({
            open: true,
            title: `Reinitialization [${instance}]`,
            content: 'Are you sure that you want to do Reinit? It will erase all node data and will download it from scratch.',
            onAgree: () => reinitNode.mutate(instance)
        })
    }
}
