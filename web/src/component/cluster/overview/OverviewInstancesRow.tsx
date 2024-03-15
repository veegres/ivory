import {Box, Radio, TableCell, TableRow, Tooltip} from "@mui/material";
import {InstanceColor, sizePretty} from "../../../app/utils";
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
import {InfoColorBox} from "../../view/box/InfoColorBox";
import {green, pink, red} from "@mui/material/colors";

const SX: SxPropsMap = {
    row: {cursor: "pointer"},
    nowrap: {whiteSpace: "nowrap", textOverflow: "ellipsis", overflow: "hidden"},
    buttons: {display: "flex", justifyContent: "flex-end", alignItems: "center"},
    data: {display: "flex", flexWrap: "wrap", columnGap: 1, rowGap: "2px", fontSize: "12px"},
}

type Props = {
    domain: string,
    instance: InstanceWeb,
    cluster: Cluster,
    candidates: Sidecar[],
}

export function OverviewInstancesRow(props: Props) {
    const {instance, domain, cluster, candidates} = props
    const {role, sidecar, database, state, lag, inInstances, pendingRestart, inCluster, scheduledRestart, scheduledSwitchover} = instance

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
            <TableCell sx={SX.nowrap}>{sidecar.host}:{sidecar.port === 0 ? "-" : sidecar.port}</TableCell>
            <TableCell sx={SX.nowrap}>{database.host}:{database.port === 0 ? "-" : database.port}</TableCell>
            <TableCell sx={SX.nowrap}>{state}</TableCell>
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
        if (role === "unknown") return null
        return (
            <Box sx={SX.data}>
                <InfoColorBox label={"Lag"} title={renderSimpleTitle("Lag", sizePretty(lag))} bgColor={lag > 100 ? red[500] : green[600]} opacity={0.9}/>
                <InfoColorBox label={"Restart"} title={renderSimpleTitle("Pending Restart", String(pendingRestart))} bgColor={pendingRestart ? green[600] : "grey"} opacity={0.9}/>
                {scheduledRestart && <InfoColorBox label={"Scheduled Restart"} title={renderScheduledRestartTitle()} bgColor={pink[900]} opacity={0.9}/>}
                {scheduledSwitchover && <InfoColorBox label={"Scheduled Switchover"} title={renderScheduledSwitchoverTitle()} bgColor={pink[900]} opacity={0.9}/>}
            </Box>
        )
    }

    function renderSimpleTitle(name: string, value: string) {
        return <Box><b>{name.toUpperCase()}:</b> {value}</Box>
    }

    function renderScheduledSwitchoverTitle() {
        if (scheduledSwitchover === undefined) return null
        return (
            <Box>
                {renderSimpleTitle("to", scheduledSwitchover.to)}
                {renderSimpleTitle("at", scheduledSwitchover.at)}
            </Box>
        )
    }

    function renderScheduledRestartTitle() {
        if (scheduledRestart === undefined) return null
        return (
            <Box>
                {renderSimpleTitle("restart", String(scheduledRestart.pendingRestart))}
                {renderSimpleTitle("at", scheduledRestart.at)}
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
