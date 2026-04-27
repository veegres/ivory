import {ErrorOutlineRounded, WarningAmberRounded} from "@mui/icons-material"
import {Box, Radio, TableCell, TableRow, Tooltip} from "@mui/material"
import {blueGrey, green, grey, pink, red} from "@mui/material/colors"

import {Cluster, Node} from "../../../../api/cluster/type"
import {Connection} from "../../../../api/node/type"
import {SxPropsMap} from "../../../../app/type"
import {
    DateTimeFormatter,
    getKeeperRequest, initialNode,
    NodeColor,
    SizeFormatter,
    SxPropsFormatter,
} from "../../../../app/utils"
import {useStoreAction} from "../../../../provider/StoreProvider"
import {InfoColorBox} from "../../../view/box/InfoColorBox"
import {MenuButton} from "../../../view/button/MenuButton"
import {HiddenScrolling} from "../../../view/scrolling/HiddenScrolling"
import {FailoverButton} from "../../../widgets/actions/FailoverButton"
import {ReinitButton} from "../../../widgets/actions/ReinitButton"
import {ReloadButton} from "../../../widgets/actions/ReloadButton"
import {RestartButton} from "../../../widgets/actions/RestartButton"
import {ScheduleButton} from "../../../widgets/actions/ScheduleButton"
import {SwitchoverButton} from "../../../widgets/actions/SwitchoverButton"

const SX: SxPropsMap = {
    row: {cursor: "pointer"},
    nowrap: {whiteSpace: "nowrap", textOverflow: "ellipsis", overflow: "hidden"},
    buttons: {display: "flex", alignItems: "center"},
    data: {display: "flex", columnGap: 1, rowGap: "2px", fontSize: "12px"},
    last: {display: "flex", justifyContent: "space-between", alignItems: "center", height: "100%"},
}

type Props = {
    name: string,
    node?: Node,
    cluster: Cluster,
    candidates: Connection[],
    checked: boolean,
    error?: boolean,
}

export function OverviewNodesRow(props: Props) {
    const {node: tmpNode, cluster, candidates, error = false, name, checked} = props
    const node = tmpNode ?? initialNode(cluster.keeperType, cluster.dbType, name)
    const {role, state, lag, pendingRestart, scheduledRestart, scheduledSwitchover, tags} = node.keeper
    const {connection, warnings} = node

    const {setNode} = useStoreAction
    const request = getKeeperRequest(cluster, connection.host, connection.keeperPort ?? 0)

    return (
        <TableRow
            sx={[SX.row, checked && SxPropsFormatter.style.bgImageSelected, error && SxPropsFormatter.style.bgImageError]}
            onClick={handleCheck}
        >
            <TableCell><Radio checked={checked} size={"small"}/></TableCell>
            <TableCell align={"center"}>{renderWarning()}</TableCell>
            <TableCell sx={{color: NodeColor[role].color}}>{role.toUpperCase()}</TableCell>
            <TableCell sx={SX.nowrap}>{node.connection.host}</TableCell>
            <TableCell sx={SX.nowrap}>{node.connection.sshPort ?? "-"}</TableCell>
            <TableCell sx={SX.nowrap}>{node.connection.keeperPort ?? "-"}</TableCell>
            <TableCell sx={SX.nowrap}>{node.connection.dbPort ?? "-"}</TableCell>
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
        if (role === "unknown") return
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
        if (role === "unknown") return
        return (
            <MenuButton>
                <ScheduleButton request={request} cluster={cluster.name} switchover={scheduledSwitchover} restart={scheduledRestart}/>
                <FailoverButton request={request} cluster={cluster.name} disabled={role === "leader"}/>
                <RestartButton request={request} cluster={cluster.name}/>
                <ReloadButton request={request} cluster={cluster.name}/>
            </MenuButton>
        )
    }

    function renderRoleButtons() {
        switch (role) {
            case "replica": return <ReinitButton request={request} cluster={cluster.name}/>
            case "leader": return <SwitchoverButton request={request} cluster={cluster.name} candidates={candidates} leaderKey={node.keeper.key}/>
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
        setNode(checked ? undefined : name)
    }
}
