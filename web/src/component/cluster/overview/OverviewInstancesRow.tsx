import {Box, Radio, TableCell, TableRow, Tooltip} from "@mui/material";
import {InstanceColor} from "../../../app/utils";
import {SxPropsMap} from "../../../type/common";
import {InstanceRequest, InstanceWeb} from "../../../type/instance";
import {WarningAmberRounded} from "@mui/icons-material";
import {useStore, useStoreAction} from "../../../provider/StoreProvider";
import {Cluster} from "../../../type/cluster";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {instanceApi} from "../../../app/api";
import {useEffect} from "react";
import {AlertButton} from "../../view/button/AlertButton";
import {MenuButton} from "../../view/button/MenuButton";

const SX: SxPropsMap = {
    row: {cursor: "pointer"},
    cell: {padding: "5px 10px", height: "50px", whiteSpace: "nowrap", textOverflow: "ellipsis", overflow: "hidden"},
    cellSmall: {padding: "5px 0", height: "50px"},
    buttons: {display: "flex", justifyContent: "flex-end", alignItems: "center"},
    data: {display: "flex", gap: 1, alignItems: "center"},
}

type Props = {
    domain: string,
    instance: InstanceWeb,
    cluster: Cluster,
}

export function OverviewInstancesRow(props: Props) {
    const {instance, domain, cluster} = props
    const {role, sidecar, database, state, lag, inInstances, pendingRestart, inCluster} = instance

    const {isInstanceActive} = useStore()
    const {setInstance} = useStoreAction()
    const checked = isInstanceActive(domain)

    const options = useMutationOptions([["instance/overview", cluster.name]])
    const switchover = useMutation({mutationFn: instanceApi.switchover, ...options})
    const reinit = useMutation({mutationFn: instanceApi.reinitialize, ...options})
    const restart = useMutation({mutationFn: instanceApi.restart, ...options})
    const reload = useMutation({mutationFn: instanceApi.reload, ...options})
    const failover = useMutation({mutationFn: instanceApi.failover, ...options})

    const request: InstanceRequest = {
        sidecar: instance.sidecar,
        credentialId: cluster.credentials.patroniId,
        certs: cluster.certs,
    }

    useEffect(handleEffectInstanceChanged, [checked, instance, setInstance])

    return (
        <TableRow sx={SX.row} onClick={handleCheck(instance, checked)}>
            <TableCell sx={SX.cell}><Radio checked={checked} size={"small"}/></TableCell>
            <TableCell sx={SX.cellSmall} align={"center"}>{renderWarning(inCluster, inInstances)}</TableCell>
            <TableCell sx={{color: InstanceColor[role]}}>{role.toUpperCase()}</TableCell>
            <TableCell sx={SX.cell} align={"center"}>
                {sidecar.host}:{sidecar.port === 0 ? "-" : sidecar.port}
            </TableCell>
            <TableCell sx={SX.cell} align={"center"}>
                {database.host}:{database.port === 0 ? "-" : database.port}
            </TableCell>
            <TableCell sx={SX.cell} align={"center"}>{state}</TableCell>
            <TableCell sx={SX.cell}>
                {renderData()}
            </TableCell>
            <TableCell sx={SX.cell} align={"right"} onClick={(e) => e.stopPropagation()}>
                <Box sx={SX.buttons}>
                    {renderButtonByRole()}
                    {renderButtons()}
                </Box>
            </TableCell>
        </TableRow>
    )

    function renderData() {
        return (
            <Box sx={SX.data}>
                {lag !== -1 && (<Box>Lag: {lag}</Box>)}
                {pendingRestart !== undefined && (<Box>Restart: {`${pendingRestart}`}</Box>)}
            </Box>
        )
    }

    function renderButtons() {
        if (role === "unknown") return null
        return (
            <MenuButton size={"small"}>
                <AlertButton
                    color={"error"}
                    size={"small"}
                    title={`Make a failover of ${instance.sidecar.host}?`}
                    content={`It will failover to current instance of postgres, that will cause some downtime 
                    and potential data loss. Usually it is recommended to use switchover, but if you don't have a
                    leader you won't be able to do switchover and here failover can help.`}
                    loading={failover.isPending}
                    disabled={instance.leader}
                    onClick={() => failover.mutate({...request, body: {candidate: instance.sidecar.host}})}
                >
                    Failover
                </AlertButton>
                <AlertButton
                    color={"primary"}
                    size={"small"}
                    title={`Make a restart of ${instance.sidecar.host}?`}
                    content={"It will restart postgres, that will cause some downtime."}
                    loading={restart.isPending}
                    onClick={() => restart.mutate(request)}
                >
                    Restart
                </AlertButton>
                <AlertButton
                    color={"primary"}
                    size={"small"}
                    title={`Make a reload of ${instance.sidecar.host}?`}
                    content={`It will reload postgres config, it doesn't have any downtime. It won't help with pending 
                    restart, because some parameters change required postgres restart.`}
                    loading={reload.isPending}
                    onClick={() => reload.mutate(request)}
                >
                    Reload
                </AlertButton>
            </MenuButton>
        )
    }

    function renderButtonByRole() {
        switch (role) {
            case "replica": return (
                <AlertButton
                    color={"info"}
                    title={`Make a reinitialization of ${instance.sidecar.host}?`}
                    content={"It will erase all node data and will download it from scratch."}
                    loading={reinit.isPending}
                    onClick={() => reinit.mutate(request)}
                >
                    Reinit
                </AlertButton>
            )
            case "leader": return (
                <AlertButton
                    color={"secondary"}
                    title={`Make a switchover of ${instance.sidecar.host}?`}
                    content={"It will change the leader of your cluster that will cause some downtime."}
                    loading={switchover.isPending}
                    onClick={() => switchover.mutate({...request, body: {leader: instance.sidecar.host}})}
                >
                    Switchover
                </AlertButton>
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
            <Box sx={SX.warning}>
                <Tooltip title={title} placement={"top"}>
                    <WarningAmberRounded color={"warning"}/>
                </Tooltip>
            </Box>
        )
    }

    function handleEffectInstanceChanged() {
        if (checked) setInstance(instance)
    }

    function handleCheck(instance: InstanceWeb, checked: boolean) {
        return () => setInstance(checked ? undefined : instance)
    }
}
