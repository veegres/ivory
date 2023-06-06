import {Box, Radio, TableCell, TableRow, Tooltip} from "@mui/material";
import {InstanceColor} from "../../../app/utils";
import {SxPropsMap} from "../../../type/common";
import {InstanceWeb} from "../../../type/instance";
import {WarningAmberRounded} from "@mui/icons-material";
import {useStore} from "../../../provider/StoreProvider";
import {Cluster} from "../../../type/cluster";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {instanceApi} from "../../../app/api";
import {AlertDialog} from "../../view/dialog/AlertDialog";
import {useEffect, useState} from "react";
import {LoadingButton} from "@mui/lab";

const SX: SxPropsMap = {
    row: {cursor: "pointer"},
    cell: {padding: "5px 10px", height: "50px", wordBreak: "break-all"},
    cellSmall: {padding: "5px 0", height: "50px"},
}

type Props = {
    domain: string,
    instance: InstanceWeb,
    cluster: Cluster,
}

export function OverviewInstancesRow(props: Props) {
    const {instance, domain, cluster} = props
    const {role, sidecar, database, state, lag, inInstances, inCluster} = instance

    const {setInstance, isInstanceActive} = useStore()
    const checked = isInstanceActive(domain)
    const [open, setOpen] = useState<"none" | "reinit" | "switch">("none")

    const options = useMutationOptions([["instance/overview", cluster.name]])
    const switchover = useMutation(instanceApi.switchover, options)
    const reinit = useMutation(instanceApi.reinitialize, options)

    // we ignore this line cause this effect uses setInstance
    // which are always changing in this function, and it causes endless recursion
    // eslint-disable-next-line
    useEffect(handleEffectInstanceChanged, [checked, instance])

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
            <TableCell sx={SX.cell} align={"center"}>{lag}</TableCell>
            <TableCell sx={SX.cell} align={"right"} onClick={(e) => e.stopPropagation()}>
                <Box display={"flex"} justifyContent={"flex-end"} alignItems={"center"}>
                    {renderButton(instance, role)}
                </Box>
            </TableCell>
        </TableRow>
    )

    function renderButton(instance: InstanceWeb, type: string) {
        switch (type) {
            case "replica": return renderReinit()
            case "leader": return renderSwitch()
            default: return null
        }
    }

    function renderReinit() {
        return (
            <>
                <AlertDialog
                    open={open === "reinit"}
                    title={`Make a reinitialization of ${instance.sidecar.host}?`}
                    content={"It will erase all node data and will download it from scratch."}
                    onAgree={() => handleReinit(instance)}
                    onClose={() => setOpen("none")}
                />
                <LoadingButton
                    color={"primary"}
                    loading={reinit.isLoading}
                    disabled={switchover.isLoading}
                    onClick={() => setOpen("reinit")}>
                    Reinit
                </LoadingButton>
            </>
        )
    }

    function renderSwitch() {
        return (
            <>
                <AlertDialog
                    open={open === "switch"}
                    title={`Make a switchover of ${instance.sidecar.host}?`}
                    content={"It will change the leader of your cluster that will cause some downtime."}
                    onAgree={() => handleSwitchover(instance)}
                    onClose={() => setOpen("none")}
                />
                <LoadingButton
                    color={"secondary"}
                    loading={switchover.isLoading}
                    onClick={() => setOpen("switch")}>
                    Switchover
                </LoadingButton>
            </>
        )
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

    function handleSwitchover(instance: InstanceWeb) {
        switchover.mutate({
            ...instance.sidecar,
            credentialId: cluster.credentials.patroniId,
            certs: cluster.certs,
            body: {leader: instance.sidecar.host},
        })
    }

    function handleReinit(instance: InstanceWeb) {
        reinit.mutate({
            ...instance.sidecar,
            credentialId: cluster.credentials.patroniId,
            certs: cluster.certs,
        })
    }
}
