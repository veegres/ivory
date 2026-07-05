import {ErrorOutlineRounded, WarningAmberRounded} from "@mui/icons-material"
import {Box, Radio, TableCell, TableRow, Tooltip} from "@mui/material"

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
    row: {cursor: "pointer"},
    nowrap: {whiteSpace: "nowrap", textOverflow: "ellipsis", overflow: "hidden"},
    buttons: {display: "flex", alignItems: "center", width: "max-content"},
    data: {display: "flex", gap: 0.5, fontSize: "12px"},
    last: {display: "flex", justifyContent: "space-between", alignItems: "center", height: "100%"},
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
        <TableRow
            sx={[SX.row, checked && SxPropsFormatter.style.bgImageSelected, error && SxPropsFormatter.style.bgImageError]}
            onClick={handleCheck}
        >
            <TableCell><Radio checked={checked} size={"small"}/></TableCell>
            <TableCell align={"center"}>{renderWarning()}</TableCell>
            <TableCell sx={{color: NodeColor[role].color}}>{role.toUpperCase()}</TableCell>
            <TableCell sx={SX.nowrap}>{config.host}</TableCell>
            <TableCell sx={SX.nowrap}>{config.keeperPort ?? "-"}</TableCell>
            <TableCell sx={SX.nowrap}>{config.dbPort ?? "-"}</TableCell>
            <TableCell sx={SX.nowrap}>{config.sshPort ?? "-"}</TableCell>
            <TableCell sx={SX.nowrap}>{state ?? "-"}</TableCell>
            <TableCell>
                <Box sx={SX.last} onClick={(e) => e.stopPropagation()}>
                    {renderData()}
                    <Box sx={SX.buttons}>{renderButtons()}</Box>
                </Box>
            </TableCell>
        </TableRow>
    )

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
            case "leader": return <KeeperSwitchoverButton request={keeperRequest} cluster={cluster.name} candidates={candidates} leaderKey={node.keeper.key}/>
            default: return
        }
    }

    function renderWarning() {
        if (warnings.length === 0) return
        return (
            <Box sx={SX.data}>
                <Tooltip title={warnings.map((s, i) => (<Box key={i}>{i+1}. {s}</Box>))} placement={"top"}>
                    <WarningAmberRounded color={"warning"}/>
                </Tooltip>
            </Box>
        )
    }

    function handleCheck() {
        setNode(checked ? undefined : nodeKey)
    }
}
