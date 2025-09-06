import {Box, Radio, TableCell, TableRow, Tooltip} from "@mui/material";
import {
    DateTimeFormatter,
    getSidecarConnection,
    InstanceColor,
    SizeFormatter,
    SxPropsFormatter
} from "../../../../app/utils";
import {InstanceWeb, Sidecar} from "../../../../api/instance/type";
import {ErrorOutlineRounded, WarningAmberRounded} from "@mui/icons-material";
import {useStore, useStoreAction} from "../../../../provider/StoreProvider";
import {Cluster} from "../../../../api/cluster/type";
import {useEffect} from "react";
import {MenuButton} from "../../../view/button/MenuButton";
import {SwitchoverButton} from "../../../widgets/actions/SwitchoverButton";
import {FailoverButton} from "../../../widgets/actions/FailoverButton";
import {RestartButton} from "../../../widgets/actions/RestartButton";
import {ReloadButton} from "../../../widgets/actions/ReloadButton";
import {ReinitButton} from "../../../widgets/actions/ReinitButton";
import {InfoColorBox} from "../../../view/box/InfoColorBox";
import {blueGrey, green, grey, pink, red} from "@mui/material/colors";
import {ScheduleButton} from "../../../widgets/actions/ScheduleButton";
import {HiddenScrolling} from "../../../view/scrolling/HiddenScrolling";
import {SxPropsMap} from "../../../../app/type";

const SX: SxPropsMap = {
    row: {cursor: "pointer"},
    nowrap: {whiteSpace: "nowrap", textOverflow: "ellipsis", overflow: "hidden"},
    buttons: {display: "flex", justifyContent: "flex-end", alignItems: "center"},
    data: {display: "flex", columnGap: 1, rowGap: "2px", fontSize: "12px"},
}

type Props = {
    domain: string,
    instance: InstanceWeb,
    cluster: Cluster,
    candidates: Sidecar[],
    error?: boolean,
}

export function OverviewInstancesRow(props: Props) {
    const {instance, domain, cluster, candidates, error = false} = props
    const {role, sidecar, database, state, lag, inInstances, pendingRestart, inCluster, scheduledRestart, scheduledSwitchover, tags} = instance

    const {isInstanceActive} = useStore()
    const {setInstance} = useStoreAction()
    const checked = isInstanceActive(domain)
    const request = getSidecarConnection(cluster, sidecar)

    useEffect(handleEffectInstanceChanged, [checked, instance, setInstance])

    const sidecarName = `${sidecar.host}:${sidecar.port === 0 ? "-" : sidecar.port}`
    const databaseName = `${database.host}:${database.port === 0 ? "-" : database.port}`

    return (
        <TableRow
            sx={[SX.row, checked && SxPropsFormatter.style.bgImageSelected, error && SxPropsFormatter.style.bgImageError]}
            onClick={handleCheck(instance, checked)}
        >
            <TableCell><Radio checked={checked} size={"small"}/></TableCell>
            <TableCell align={"center"}>{renderWarning(inCluster, inInstances)}</TableCell>
            <TableCell sx={{color: InstanceColor[role]}}>{role.toUpperCase()}</TableCell>
            <TableCell sx={SX.nowrap}>
                <Tooltip title={sidecarName} placement={"top-start"}>
                    <Box sx={SX.nowrap}>{sidecarName}</Box>
                </Tooltip>
            </TableCell>
            <TableCell sx={SX.nowrap}>
                <Tooltip title={databaseName} placement={"top-start"}>
                    <Box sx={SX.nowrap}>{databaseName}</Box>
                </Tooltip>
            </TableCell>
            <TableCell sx={SX.nowrap}>{state}</TableCell>
            <TableCell onClick={(e) => e.stopPropagation()}>{renderData()}</TableCell>
            <TableCell align={"right"} onClick={(e) => e.stopPropagation()}>
                <Box sx={SX.buttons}>{renderButtons()}</Box>
            </TableCell>
        </TableRow>
    )

    function renderButtons() {
        if (error) return (
            <Tooltip title={"This instance is not in this cluster anymore. It was before that is why you see it. Just uncheck it."} placement={"top"}>
                <ErrorOutlineRounded color={"error"}/>
            </Tooltip>
        )

        return (
            <>
                {renderRoleButtons()}
                {renderMenuButtons()}
            </>
        )
    }

    function renderData() {
        if (role === "unknown") return null
        return (
            <HiddenScrolling arrowWidth={"20px"} arrowHeight={"25px"}>
                <Box sx={SX.data}>
                    <InfoColorBox label={"Restart"} title={renderSimpleTitle("Pending Restart", String(pendingRestart))} bgColor={pendingRestart ? green[600] : grey[600]} opacity={0.9}/>
                    {role === "replica" && <InfoColorBox label={"Lag"} title={renderSimpleTitle("Lag", SizeFormatter.pretty(lag))} bgColor={lag > 100 ? red[500] : grey[600]} opacity={0.9}/>}
                    {scheduledRestart && <InfoColorBox label={"Scheduled Restart"} title={renderScheduledRestartTitle()} bgColor={pink[900]} opacity={0.9}/>}
                    {scheduledSwitchover && <InfoColorBox label={"Scheduled Switchover"} title={renderScheduledSwitchoverTitle()} bgColor={pink[900]} opacity={0.9}/>}
                    {renderTags()}
                </Box>
            </HiddenScrolling>
        )
    }

    function renderTags() {
        if (!tags) return
        return Object.entries(tags).map(([key, value]) => (
            <InfoColorBox key={key} label={`${key}: ${value}`} bgColor={blueGrey[600]} opacity={0.9}/>
        ))
    }

    function renderSimpleTitle(name: string, value: string) {
        return <Box><b>{name.toUpperCase()}:</b> {value}</Box>
    }

    function renderScheduledSwitchoverTitle() {
        if (scheduledSwitchover === undefined) return null
        return (
            <Box>
                {renderSimpleTitle("to", scheduledSwitchover.to)}
                {renderSimpleTitle("at", DateTimeFormatter.utc(scheduledSwitchover.at))}
            </Box>
        )
    }

    function renderScheduledRestartTitle() {
        if (scheduledRestart === undefined) return null
        return (
            <Box>
                {renderSimpleTitle("pending restart", String(scheduledRestart.pendingRestart))}
                {renderSimpleTitle("at", DateTimeFormatter.utc(scheduledRestart.at))}
            </Box>
        )
    }

    function renderMenuButtons() {
        if (role === "unknown") return null
        return (
            <MenuButton>
                <ScheduleButton request={request} cluster={cluster.name} switchover={scheduledSwitchover} restart={scheduledRestart}/>
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
