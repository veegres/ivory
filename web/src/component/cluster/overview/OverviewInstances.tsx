import {Box, Button, Radio, Table, TableBody, TableCell, TableHead, TableRow, Tooltip} from "@mui/material";
import {useIsFetching, useMutation} from "@tanstack/react-query";
import {instanceApi} from "../../../app/api";
import {useMemo, useState} from "react";
import {TableCellLoader} from "../../view/TableCellLoader";
import {getDomain, InstanceColor} from "../../../app/utils";
import {AlertDialog} from "../../view/AlertDialog";
import {useStore} from "../../../provider/StoreProvider";
import {TabProps} from "./Overview";
import {Warning} from "@mui/icons-material";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {SxPropsMap} from "../../../type/common";
import {ActiveInstance} from "../../../type/instance";

const SX: SxPropsMap = {
    table: {"tr:last-child td": {border: 0}},
    row: {cursor: "pointer"},
    cell: {padding: "5px 10px", height: "50px"},
    actionCell: {width: "58px"},
    warningCell: {width: "40px"},
    roleCell: {width: "110px"},
    buttonCell: {width: "160px"},
}

type AlertDialogState = { open: boolean, title: string, content: string, onAgree: () => void }
const initAlertDialog = {
    open: false, title: "", content: "", onAgree: () => void 0
}

export function OverviewInstances({info}: TabProps) {
    const {defaultInstance, cluster, combinedInstanceMap} = info
    const [alertDialog, setAlertDialog] = useState<AlertDialogState>(initAlertDialog)
    const {setInstance, store: {activeInstance}} = useStore()
    const queryKey = useMemo(
        () => ["instance/overview", cluster.name, defaultInstance.sidecar.host, defaultInstance.sidecar.port],
        [defaultInstance.sidecar, cluster.name]
    )
    const instanceMapFetching = useIsFetching(queryKey)
    const options = useMutationOptions([queryKey])
    const switchover = useMutation(instanceApi.switchover, options)
    const reinit = useMutation(instanceApi.reinitialize, options)

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
                        <TableCellLoader
                            sx={SX.buttonCell}
                            isFetching={(instanceMapFetching > 0 || switchover.isLoading || reinit.isLoading)}
                        />
                    </TableRow>
                </TableHead>
                <TableBody>
                    {renderContent()}
                </TableBody>
            </Table>
        </>
    )

    function renderContent() {
        return Object.entries(combinedInstanceMap).map(([key, element]) => {
            const {role, sidecar, database, state, lag, inInstances, inCluster} = element
            const isChecked = activeInstance ? getDomain(activeInstance.sidecar) === key : false
            const instance = {cluster: cluster.name, sidecar: sidecar, database: database}

            return (
                <TableRow sx={SX.row} key={key} onClick={handleCheck(instance, isChecked)}>
                    <TableCell sx={SX.cell}><Radio checked={isChecked} size={"small"}/></TableCell>
                    <TableCell sx={SX.cell} align={"center"}>{renderWarning(inCluster, inInstances)}</TableCell>
                    <TableCell sx={{color: InstanceColor[role]}}>{role.toUpperCase()}</TableCell>
                    <TableCell sx={SX.cell} align={"center"}>
                        {sidecar.host}:{sidecar.port === 0 ? "-" : sidecar.port}
                    </TableCell>
                    <TableCell sx={SX.cell} align={"center"}>
                        {database.host}:{database.port === 0 ? "-" : database.port}
                    </TableCell>
                    <TableCell sx={SX.cell} align={"center"}>{state}</TableCell>
                    <TableCell sx={SX.cell} align={"center"}>{lag}</TableCell>
                    <TableCell sx={SX.cell} align={"right"} onClick={(e) => e.stopPropagation()}>
                        <Box display={"flex"} justifyContent={"flex-end"} alignItems={"center"}>
                            {renderButton(instance, role)}
                        </Box>
                    </TableCell>
                </TableRow>
            )
        })
    }

    function renderButton(instance: ActiveInstance, type: string) {
        switch (type) {
            case "replica":
                return (
                    <Button
                        color={"primary"}
                        disabled={reinit.isLoading || switchover.isLoading}
                        onClick={() => handleReinit(instance)}>
                        Reinit
                    </Button>
                )
            case "leader":
                return (
                    <Button
                        color={"secondary"}
                        disabled={switchover.isLoading}
                        onClick={() => handleSwitchover(instance)}>
                        Switchover
                    </Button>
                )
            default:
                return null
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
                    <Warning color={"warning"} fontSize={"small"}/>
                </Tooltip>
            </Box>
        )
    }

    function handleCheck(instance: ActiveInstance, checked: boolean) {
        return () => setInstance(checked ? undefined : instance)
    }

    function handleSwitchover(instance: ActiveInstance) {
        setAlertDialog({
            open: true,
            title: `Switchover [${instance.sidecar.host}]`,
            content: "Are you sure that you want to do Switchover? It will change the leader of your cluster.",
            onAgree: () => switchover.mutate({
                cluster: instance.cluster,
                host: instance.sidecar.host,
                port: instance.sidecar.port,
                body: {leader: instance.sidecar.host},
            })
        })
    }

    function handleReinit(instance: ActiveInstance) {
        setAlertDialog({
            open: true,
            title: `Reinitialization [${instance.sidecar.host}]`,
            content: "Are you sure that you want to do Reinit? It will erase all node data and will download it from scratch.",
            onAgree: () => reinit.mutate({
                cluster: instance.cluster,
                host: instance.sidecar.host,
                port: instance.sidecar.port,
            })
        })
    }
}
