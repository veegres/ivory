import {ErrorOutlineRounded, WarningAmberRounded} from "@mui/icons-material"
import {Box, Tooltip} from "@mui/material"

import {Cluster, Node, NodeConfig} from "../../../../features/cluster/api/ClusterType"
import {KeeperFailoverButton} from "../../../../features/node/component/keeper/KeeperFailoverButton"
import {KeeperReinitButton} from "../../../../features/node/component/keeper/KeeperReinitButton"
import {KeeperReloadButton} from "../../../../features/node/component/keeper/KeeperReloadButton"
import {KeeperRestartButton} from "../../../../features/node/component/keeper/KeeperRestartButton"
import {KeeperScheduleButton} from "../../../../features/node/component/keeper/KeeperScheduleButton"
import {KeeperSwitchoverButton} from "../../../../features/node/component/keeper/KeeperSwitchoverButton"
import {InfoColorBox} from "../../../../shared/component/box/InfoColorBox"
import {MenuButton} from "../../../../shared/component/button/MenuButton"
import {HiddenScrolling} from "../../../../shared/component/scrolling/HiddenScrolling"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {DateTimeFormatter, getKeeperOneRequest, NodeColor, SizeFormatter, SxPropsFormatter} from "../../../../shared/helper/HelperUtils"
import {useStoreAction} from "../../../../shared/provider/StoreProvider"

const SX: SxPropsMap = {
    row: {
        display: "flex", flexWrap: "wrap", justifyContent: "space-between", alignItems: "center",
        gap: 0.5, padding: "2px 5px", minHeight: "38px", cursor: "pointer",
        borderBottom: 1, borderColor: "divider",
        "&:last-child": {borderBottom: 0}, borderLeft: "3px solid transparent",
    },
    checked: {borderLeftColor: "primary.main"},
    group: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 0.5, flexGrow: 1},
    // NOTE: the warning slot is always rendered (empty when there are no
    // warnings) so rows keep the same layout either way
    warning: {flex: "0 0 28px", display: "flex", justifyContent: "center"},
    // NOTE: the host shrinks with ellipsis when long but doesn't grow when
    // short, so the leftover space goes to the data/actions segment; the min
    // width and right padding keep the segment a bit wider than the name
    host: {flex: "0 1 auto", minWidth: "120px", maxWidth: "100%", paddingRight: "16px"},
    hostValue: {whiteSpace: "nowrap", textOverflow: "ellipsis", overflow: "hidden", maxWidth: "100%"},
    info: {display: "flex", flexDirection: "column", alignItems: "flex-start", whiteSpace: "nowrap", lineHeight: 1.2, minWidth: 0, padding: "0 6px"},
    infoLabel: {fontSize: "8px", textTransform: "uppercase", color: "text.secondary"},
    buttons: {display: "flex", alignItems: "center", width: "max-content"},
    data: {display: "flex", gap: 0.5, fontSize: "12px"},
    // NOTE: the grow factor dwarfs the groups' 1, so on a shared line the
    // status/actions segment takes practically all the free space while the
    // groups stay at content size; a group only stretches when it wraps to a
    // line of its own
    last: {
        flex: "1000 1 260px", minWidth: 0, display: "flex",
        justifyContent: "space-between", alignItems: "center",
    },
    title: {fontFamily: "monospace", textTransform: "uppercase"},
}

type Props = {
    node: Node,
    nodeKey: string,
    cluster: Cluster,
    candidates: NodeConfig[],
    checked: boolean,
    error?: boolean,
}

export function OverviewNodesRow(props: Props) {
    const {node, cluster, candidates, error = false, nodeKey, checked} = props
    const {config, warnings, keeper} = node
    const {role, state, lag, pendingRestart, scheduledRestart, scheduledSwitchover, tags} = keeper

    const {setNode} = useStoreAction
    const keeperRequest = getKeeperOneRequest(cluster, config.host, config.keeperPort)

    return (
        <Box
            sx={[SX.row, checked && SX.checked, checked && SxPropsFormatter.style.bgImageSelected, error && SxPropsFormatter.style.bgImageError]}
            onClick={handleCheck}
        >
            <Box sx={SX.group}>
                {renderHost()}
                {renderRole()}
            </Box>
            <Box sx={SX.group}>
                {renderInfo("keeper", config.keeperPort, "70px")}
                {renderInfo("db", config.dbPort, "70px")}
                {renderInfo("ssh", config.sshPort, "70px")}
            </Box>
            <Box sx={SX.group}>
                {renderState()}
                <Box sx={SX.warning}>{renderWarning()}</Box>
            </Box>
            <Box sx={SX.last} onClick={(e) => e.stopPropagation()}>
                {renderData()}
                <Box sx={SX.buttons}>{renderButtons()}</Box>
            </Box>
        </Box>
    )

    function renderHost() {
        return (
            <Box sx={[SX.info, SX.host]}>
                <Box sx={SX.infoLabel}>host</Box>
                <Box sx={SX.hostValue}>{config.host}</Box>
            </Box>
        )
    }

    // NOTE: the fixed min width keeps LEADER and REPLICA the same size
    function renderRole() {
        return (
            <Box sx={[SX.info, {minWidth: "100px"}]}>
                <Box sx={SX.infoLabel}>role</Box>
                <Box sx={{color: NodeColor[role].color}}>{role.toUpperCase()}</Box>
            </Box>
        )
    }

    function renderState() {
        return (
            <Box sx={[SX.info, {minWidth: "100px", fontSize: "12px"}]}>
                <Box sx={SX.infoLabel}>state</Box>
                <InfoColorBox label={state ?? "unknown"} color={getStateColor()}/>
            </Box>
        )
    }

    // NOTE: the min width applies at every breakpoint so a port segment keeps
    // the same size whether it has a value or shows "-"
    function renderInfo(label: string, value?: string | number, width?: string) {
        return (
            <Box sx={[SX.info, {minWidth: width}]}>
                <Box sx={SX.infoLabel}>{label}</Box>
                <Box>{value ?? "-"}</Box>
            </Box>
        )
    }

    function renderButtons() {
        if (error) return (
            <Tooltip title={"This node is not in this cluster anymore. It was before that is why you see it. Just uncheck it."} placement={"top"}>
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
        if (role === "unknown") return <Box/>
        return (
            <HiddenScrolling arrowWidth={"20px"} arrowHeight={"25px"}>
                <Box sx={SX.data}>
                    {pendingRestart && <InfoColorBox label={"Restart"} dot={true} title={renderSimpleTitle("Pending Restart", String(pendingRestart))} color={"warning"}/>}
                    {role === "replica" && <InfoColorBox label={"Lag"} dot={true} title={renderSimpleTitle("Lag", SizeFormatter.pretty(lag))} color={lag > 100 ? "error" : "default"}/>}
                    {scheduledRestart && <InfoColorBox label={"Scheduled Restart"} title={renderScheduledRestartTitle()} color={"secondary"}/>}
                    {scheduledSwitchover && <InfoColorBox label={"Scheduled Switchover"} title={renderScheduledSwitchoverTitle()} color={"secondary"}/>}
                    {renderTags()}
                </Box>
            </HiddenScrolling>
        )
    }

    function renderTags() {
        if (!tags) return
        return Object.entries(tags).map(([key, value]) => (
            <InfoColorBox key={key} label={`${key}: ${value}`} color={"info"}/>
        ))
    }

    function renderSimpleTitle(name: string, value: string) {
        return <Box sx={SX.title}>{name.toUpperCase()}: <b>{value}</b></Box>
    }

    function renderScheduledSwitchoverTitle() {
        if (scheduledSwitchover === undefined) return
        return (
            <Box>
                {renderSimpleTitle("to", scheduledSwitchover.to)}
                {renderSimpleTitle("at", DateTimeFormatter.utc(scheduledSwitchover.at))}
            </Box>
        )
    }

    function renderScheduledRestartTitle() {
        if (scheduledRestart === undefined) return
        return (
            <Box>
                {renderSimpleTitle("pending restart", String(scheduledRestart.pendingRestart))}
                {renderSimpleTitle("at", DateTimeFormatter.utc(scheduledRestart.at))}
            </Box>
        )
    }

    function renderMenuButtons() {
        if (role === "unknown" || !keeperRequest) return
        return (
            <MenuButton size={27}>
                <KeeperScheduleButton request={keeperRequest} cluster={cluster.name} switchover={scheduledSwitchover} restart={scheduledRestart}/>
                <KeeperFailoverButton request={keeperRequest} cluster={cluster.name} role={role}/>
                <KeeperRestartButton request={keeperRequest} cluster={cluster.name}/>
                <KeeperReloadButton request={keeperRequest} cluster={cluster.name}/>
            </MenuButton>
        )
    }

    function renderRoleButtons() {
        if (!keeperRequest) return
        switch (role) {
            case "replica": return <KeeperReinitButton request={keeperRequest} cluster={cluster.name}/>
            case "leader": return <KeeperSwitchoverButton request={keeperRequest} cluster={cluster.name} candidates={candidates.map(c => c.host)} leaderKey={node.keeper.key}/>
            default: return
        }
    }

    function renderWarning() {
        if (warnings.length === 0) return
        return (
            <Tooltip title={warnings.map((s, i) => (<Box key={i}>{i+1}. {s}</Box>))} placement={"top"}>
                <WarningAmberRounded color={"warning"}/>
            </Tooltip>
        )
    }

    function handleCheck() {
        setNode(checked ? undefined : nodeKey)
    }

    function getStateColor(): "success" | "warning" | "error" | "default" {
        switch (state) {
            case "running":
                return "success"
            case "starting":
            case "restarting":
            case "stopping":
                return "warning"
            case "stopped":
            case "failed":
            case "unreachable":
                return "error"
            default:
                return "default"
        }
    }
}
