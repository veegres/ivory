import {Box, Button, Radio, Table, TableBody, TableCell, TableHead, TableRow, Tooltip} from "@mui/material";
import {useMutation, useQuery} from "@tanstack/react-query";
import {nodeApi} from "../../../app/api";
import React, {useState} from "react";
import {TableCellLoader} from "../../view/TableCellLoader";
import {InstanceColor} from "../../../app/utils";
import {AlertDialog} from "../../view/AlertDialog";
import {useStore} from "../../../provider/StoreProvider";
import {TabProps} from "./Overview";
import {Warning} from "@mui/icons-material";
import {useMutationOptions} from "../../../hook/QueryCustom";

const SX = {
    table: {'tr:last-child td': {border: 0}},
    row: {cursor: "pointer"},
    cell: {padding: "5px 10px", height: "50px"},
    actionCell: {width: "58px"},
    warningCell: {width: "40px"},
    roleCell: {width: "110px"},
    buttonCell: {width: "160px"},
}

type AlertDialogState = {open: boolean, title: string, content: string, onAgree: () => void}
const initAlertDialog = {open: false, title: '', content: '', onAgree: () => {}}

export function OverviewInstances({info}: TabProps) {
    const { instance, cluster, instances } = info
    const [alertDialog, setAlertDialog] = useState<AlertDialogState>(initAlertDialog)
    const { setInstance, store: { activeInstance } } = useStore()

    const instanceMap = useQuery(
        ["node/cluster", cluster.name, instance.api_domain],
        () => nodeApi.cluster(instance.api_domain),
        { retry: 0 }
    )
    const switchoverNodeOptions = useMutationOptions(["node/cluster", cluster.name, instance.api_domain])
    const switchoverNode = useMutation(nodeApi.switchover, switchoverNodeOptions)
    const reinitNodeOptions = useMutationOptions(["node/cluster", cluster.name, instance.api_domain])
    const reinitNode = useMutation(nodeApi.reinitialize, reinitNodeOptions)

    return (
        <>
            <AlertDialog {...alertDialog} onClose={() => setAlertDialog({...alertDialog, open: false})}/>
            <Table size={"small"} sx={SX.table}>
                <TableHead>
                    <TableRow>
                        <TableCell sx={SX.actionCell}/>
                        <TableCell sx={SX.warningCell}/>
                        <TableCell sx={SX.roleCell}>Role</TableCell>
                        <TableCell align={"center"}>Patroni</TableCell>
                        <TableCell align={"center"}>Postgres</TableCell>
                        <TableCell align={"center"}>State</TableCell>
                        <TableCell align={"center"}>Lag</TableCell>
                        <TableCellLoader sx={SX.buttonCell} isFetching={(instanceMap.isFetching || switchoverNode.isLoading || reinitNode.isLoading)}/>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {renderContent()}
                </TableBody>
            </Table>
        </>
    )

    function renderContent() {
        return Object.entries(instances).map(([key, element]) => {
            const { role, port, host, state, lag, inInstances, inCluster } = element
            const isChecked = activeInstance === key

            return (
                <TableRow sx={SX.row} key={key} onClick={handleCheck(isChecked, key)}>
                    <TableCell sx={SX.cell}><Radio checked={isChecked} size={"small"}/></TableCell>
                    <TableCell sx={SX.cell} align={"center"}>{renderWarning(inCluster, inInstances)}</TableCell>
                    <TableCell sx={{color: InstanceColor[role]}}>{role.toUpperCase()}</TableCell>
                    <TableCell sx={SX.cell} align={"center"}>{key}</TableCell>
                    <TableCell sx={SX.cell} align={"center"}>{host}:{port ? port : "-"}</TableCell>
                    <TableCell sx={SX.cell} align={"center"}>{state}</TableCell>
                    <TableCell sx={SX.cell} align={"center"}>{lag}</TableCell>
                    <TableCell sx={SX.cell} align={"right"} onClick={(e) => e.stopPropagation()}>
                        <Box display={"flex"} justifyContent={"flex-end"} alignItems={"center"}>
                            {renderButton(key, host, role)}
                        </Box>
                    </TableCell>
                </TableRow>
            )
        })
    }

    function renderButton(instance: string, host: string, type: string) {
        switch (type) {
            case "replica": return (
                <Button
                    color={"primary"}
                    disabled={reinitNode.isLoading || switchoverNode.isLoading}
                    onClick={() => handleReinit(instance)}>
                    Reinit
                </Button>
            )
            case "leader": return (
                <Button
                    color={"secondary"}
                    disabled={switchoverNode.isLoading}
                    onClick={() => handleSwitchover(instance, host)}>
                    Switchover
                </Button>
            )
            default: return null
        }

    }

    function renderWarning(inCluster: boolean, inInstances: boolean) {
        if (inCluster && inInstances) return null

        let title = "Instance is not in cluster list or cluster itself!"
        if (inCluster && !inInstances) title = "Instance is not in your cluster list!"
        if (!inCluster && inInstances) title = "Instance is not in cluster!"

        return (
            <Box display={"flex"}>
                <Tooltip title={title} placement={"top"}>
                    <Warning color={"warning"} fontSize={"small"} />
                </Tooltip>
            </Box>
        )
    }

    function handleCheck(checked: boolean, domain: string) {
        return () => setInstance(checked ? '' : domain)
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
