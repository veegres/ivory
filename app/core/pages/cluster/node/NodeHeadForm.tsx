import {Box, TextField} from "@mui/material"
import {useState} from "react"

import {Node, NodeConfig} from "../../../../features/cluster/api/ClusterType"
import {KeeperState, KeeperStatus} from "../../../../features/node/api/NodeType"
import {InfoColorBox} from "../../../../shared/component/box/InfoColorBox"
import {InfoStatusItem, InfoStatusList} from "../../../../shared/component/box/InfoStatusList"
import {CancelIconButton, EditIconButton, SaveIconButton} from "../../../../shared/component/button/IconButtons"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {SizeFormatter} from "../../../../shared/helper/HelperUtils"

const SX: SxPropsMap = {
    box: {display: "flex", width: "100%", padding: "6px 8px", gap: 1.5, alignItems: "center", flex: "11 1 700px"},
    body: {display: "flex", flexDirection: "column", width: "100%", gap: 1.5},
    container: {
        display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(150px, 1fr))", gap: 1.5,
        alignItems: "center", flex: 1,
    },
    actions: {display: "flex", flexDirection: "column", height: "100%"},
    tags: {display: "flex", gap: 0.5},
}

type Props = {
    node: Node,
    loading?: boolean,
    onUpdate?: (config: NodeConfig) => void,
}

export function NodeHeadForm(props: Props) {
    const {node, onUpdate, loading} = props
    const [edit, setEdit] = useState(false)
    const [config, setConfig] = useState(node.config)

    const stateStr = node.keeper.state ?? "unknown"
    const isPaused = node.keeper.status === KeeperStatus.Paused
    const stateColor = isPaused ? "warning" : getStateColor(stateStr)

    const isReplica = node.keeper.role === "replica"
    const lagValue = isReplica ? SizeFormatter.pretty(node.keeper.lag) : "N/A"
    const lagColor = !isReplica ? "default" : (node.keeper.lag > 100 * 1024 * 1024 ? "error" : (node.keeper.lag > 0 ? "warning" : "success"))

    const pending = node.keeper.pendingRestart
    const pendingText = pending ? "Pending" : "No"
    const pendingColor = pending ? "warning" : "success"

    const schedRestart = node.keeper.scheduledRestart
    const restartText = schedRestart ? "Scheduled" : "None"
    const restartColor = schedRestart ? "secondary" : "default"

    const schedSwitch = node.keeper.scheduledSwitchover
    const switchText = schedSwitch ? `To ${schedSwitch.to}` : "None"
    const switchColor = schedSwitch ? "secondary" : "default"

    return (
        <Box sx={SX.box}>
            <Box sx={SX.body}>
                <Box sx={SX.container}>
                    {renderItem("Host", node.config.host, "host", "text")}
                    {renderItem("Keeper", node.config.keeperPort, "keeperPort", "number")}
                    {renderItem("Database", node.config.dbPort, "dbPort", "number")}
                    {renderItem("SSH", node.config.sshPort, "sshPort", "number")}
                </Box>
                <InfoStatusList>
                    {renderStatusItem("State", stateStr, stateColor)}
                    {renderStatusItem("Lag", lagValue, lagColor)}
                    {renderStatusItem("Pending Restart", pendingText, pendingColor)}
                    {renderStatusItem("Scheduled Restart", restartText, restartColor)}
                    {renderStatusItem("Scheduled Switchover", switchText, switchColor)}
                    {renderTagsItem()}
                </InfoStatusList>
            </Box>
            {renderActions()}
        </Box>
    )

    function renderItem(label: string, value: string | number | undefined, field: keyof NodeConfig, type: "text" | "number") {
        return (
            <TextField
                fullWidth
                size={"small"}
                label={label}
                type={type}
                value={edit ? (config[field] ?? "") : (value ?? "")}
                disabled={!edit || loading}
                onChange={e => {
                    const val = e.target.value
                    const parsed = type === "number" ? (val === "" ? undefined : parseInt(val)) : val
                    setConfig({...config, [field]: parsed})
                }}
            />
        )
    }

    function renderStatusItem(label: string, text: string, color: "success" | "warning" | "error" | "info" | "secondary" | "default") {
        return (
            <InfoStatusItem label={label}>
                <InfoColorBox label={text} color={color} dot={true}/>
            </InfoStatusItem>
        )
    }

    function renderTagsItem() {
        const tagsEntries = node.keeper.tags ? Object.entries(node.keeper.tags) : []

        return (
            <InfoStatusItem label={"Tags"}>
                <Box sx={SX.tags}>
                    {tagsEntries.length === 0 ? (
                        <InfoColorBox label={"None"} dot={true}/>
                    ) : tagsEntries.map(([key, value]) => (
                        <InfoColorBox key={key} label={`${key}: ${value}`} color={"info"}/>
                    ))}
                </Box>
            </InfoStatusItem>
        )
    }

    function renderActions() {
        if (!onUpdate) return
        return (
            <Box sx={SX.actions}>
                {edit ? (
                    <>
                        <CancelIconButton size={32} onClick={handleCancel} disabled={loading}/>
                        <SaveIconButton size={32} onClick={handleSave} loading={loading}/>
                    </>
                ) : (
                    <EditIconButton size={32} onClick={handleEdit}/>
                )}
            </Box>
        )
    }

    function handleSave() {
        onUpdate?.(config)
        setEdit(false)
    }

    function handleCancel() {
        setEdit(false)
        setConfig(node.config)
    }

    function handleEdit() {
        setEdit(true)
        setConfig(node.config)
    }

    function getStateColor(state: KeeperState): "success" | "warning" | "error" | "default" {
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
