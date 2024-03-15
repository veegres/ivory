import {Box, Radio, TableCell, TableRow, Tooltip} from "@mui/material";
import {InstanceColor} from "../../../app/utils";
import {Sidecar, SxPropsMap} from "../../../type/common";
import {InstanceRequest, InstanceWeb} from "../../../type/instance";
import {WarningAmberRounded} from "@mui/icons-material";
import {useStore, useStoreAction} from "../../../provider/StoreProvider";
import {Cluster} from "../../../type/cluster";
import {useEffect} from "react";
import {MenuButton} from "../../view/button/MenuButton";
import {SwitchoverButton} from "../../shared/actions/SwitchoverButton";
import {FailoverButton} from "../../shared/actions/FailoverButton";
import {RestartButton} from "../../shared/actions/RestartButton";
import {ReloadButton} from "../../shared/actions/ReloadButton";
import {ReinitButton} from "../../shared/actions/ReinitButton";

const SX: SxPropsMap = {
    row: {cursor: "pointer"},
    nowrap: {whiteSpace: "nowrap", textOverflow: "ellipsis", overflow: "hidden"},
    buttons: {display: "flex", justifyContent: "flex-end", alignItems: "center"},
    data: {display: "flex", gap: 1, alignItems: "center"},
}

type Props = {
    domain: string,
    instance: InstanceWeb,
    cluster: Cluster,
    candidates: Sidecar[],
}

export function OverviewInstancesRow(props: Props) {
    const {instance, domain, cluster, candidates} = props
    const {role, sidecar, database, state, lag, inInstances, pendingRestart, inCluster} = instance

    const {isInstanceActive} = useStore()
    const {setInstance} = useStoreAction()
    const checked = isInstanceActive(domain)

    const request: InstanceRequest = {
        sidecar: instance.sidecar,
        credentialId: cluster.credentials.patroniId,
        certs: cluster.certs,
    }

    useEffect(handleEffectInstanceChanged, [checked, instance, setInstance])

    return (
        <TableRow sx={SX.row} onClick={handleCheck(instance, checked)}>
            <TableCell><Radio checked={checked} size={"small"}/></TableCell>
            <TableCell align={"center"}>{renderWarning(inCluster, inInstances)}</TableCell>
            <TableCell sx={{color: InstanceColor[role]}}>{role.toUpperCase()}</TableCell>
            <TableCell sx={SX.nowrap} align={"center"}>{sidecar.host}:{sidecar.port === 0 ? "-" : sidecar.port}</TableCell>
            <TableCell sx={SX.nowrap} align={"center"}>{database.host}:{database.port === 0 ? "-" : database.port}</TableCell>
            <TableCell sx={SX.nowrap} align={"center"}>{state}</TableCell>
            <TableCell>{renderData()}</TableCell>
            <TableCell align={"right"} onClick={(e) => e.stopPropagation()}>
                <Box sx={SX.buttons}>
                    {renderRoleButtons()}
                    {renderMenuButtons()}
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

    function renderMenuButtons() {
        if (role === "unknown") return null
        return (
            <MenuButton size={"small"}>
                <FailoverButton request={request} cluster={cluster.name} disabled={instance.leader}/>
                <RestartButton request={request} cluster={cluster.name}/>
                <ReloadButton request={request} cluster={cluster.name}/>
            </MenuButton>
        )
    }

    function renderRoleButtons() {
        switch (role) {
            case "replica": return <ReinitButton request={request} cluster={cluster.name}/>
            case "leader": return <SwitchoverButton request={request} cluster={cluster.name} candidates={candidates}/>
            default: return null
        }
    }

    function renderWarning(inCluster: boolean, inInstances: boolean) {
        if (inCluster && inInstances) return null

        let title = "Instance is not in cluster list or cluster itself!"
        if (inCluster && !inInstances) title = "Instance is not in your cluster list!"
        if (!inCluster && inInstances) title = "Instance is not in cluster!"

        return (
            <Box sx={SX.data}>
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
