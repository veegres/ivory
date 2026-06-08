import {Box, TextField} from "@mui/material"
import {useState} from "react"

import {Node, NodeConfig} from "../../../../features/cluster/type"
import {CancelIconButton, EditIconButton, SaveIconButton} from "../../../../shared/component/button/IconButtons"
import {SxPropsMap} from "../../../../shared/helper/type"
import {SizeFormatter} from "../../../../shared/helper/utils"

const SX: SxPropsMap = {
    box: {display: "flex", width: "100%", padding: "6px 8px", gap: 1.5, alignItems: "center"},
    body: {display: "flex", flexDirection: "column", width: "100%", gap: 1.5},
    container: {
        display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(150px, 1fr))", gap: 1.5,
        alignItems: "center", flex: 1,
    },
    statusContainer: {
        display: "flex", flexWrap: "wrap", gap: "8px 20px", alignItems: "center", width: "100%",
    },
    actions: {display: "flex", flexDirection: "column", height: "100%"},
    item: {display: "flex", flexDirection: "column", gap: 0.25},
    title: {
        color: "text.disabled", fontWeight: "bold", fontSize: "0.7rem", textTransform: "uppercase",
        letterSpacing: "0.05em", whiteSpace: "nowrap",
    },
    badge: {
        display: "inline-flex", alignItems: "center", gap: 0.75, px: 1, borderRadius: "6px",
        fontWeight: 600, fontSize: "0.75rem", width: "fit-content", whiteSpace: "nowrap",
    },
    dot: {width: 5, height: 5, borderRadius: "50%"},
}

type Props = {
    node: Node,
    loading?: boolean,
    onUpdate?: (config: NodeConfig) => void,
}

export function NodeInfoForm(props: Props) {
    const {node, onUpdate, loading} = props

    const [edit, setEdit] = useState(false)
    const [config, setConfig] = useState<NodeConfig>(node.config)

    const stateStr = node.keeper.state ?? "unknown"
    const isRunning = stateStr.toLowerCase() === "running"
    const isPaused = stateStr.toLowerCase() === "paused"
    const isUnknown = stateStr.toLowerCase() === "unknown"
    const stateColor = isRunning ? "success" : (isPaused ? "warning" : (isUnknown ? "default" : "error"))

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
                    {renderItem("SSH", node.config.sshPort, "sshPort", "number")}
                    {renderItem("Keeper", node.config.keeperPort, "keeperPort", "number")}
                    {renderItem("Database", node.config.dbPort, "dbPort", "number")}
                </Box>
                <Box sx={SX.statusContainer}>
                    {renderStatusItem("State", stateStr, stateColor)}
                    {renderStatusItem("Lag", lagValue, lagColor)}
                    {renderStatusItem("Pending Restart", pendingText, pendingColor)}
                    {renderStatusItem("Scheduled Restart", restartText, restartColor)}
                    {renderStatusItem("Scheduled Switchover", switchText, switchColor)}
                    {renderTagsItem()}
                </Box>
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

    function renderStatusItem(label: string, text: string, colorType: "success" | "warning" | "error" | "info" | "secondary" | "default") {
        return (
            <Box sx={SX.item}>
                <Box sx={SX.title}>{label}</Box>
                {renderBadge(text, colorType)}
            </Box>
        )
    }

    function renderTagsItem() {
        const tagsEntries = node.keeper.tags ? Object.entries(node.keeper.tags) : []

        return (
            <Box sx={SX.item}>
                <Box sx={SX.title}>Tags</Box>
                {tagsEntries.length === 0 ? (
                    renderBadge("None", "default")
                ) : (
                    <Box sx={{display: "flex", flexWrap: "wrap", gap: 0.5}}>
                        {tagsEntries.map(([key, value]) => (
                            <Box
                                key={key}
                                sx={{
                                    px: 1,
                                    py: 0.25,
                                    borderRadius: "8px",
                                    bgcolor: "action.hover",
                                    fontSize: "0.75rem",
                                    fontWeight: 500,
                                    color: "text.secondary",
                                    border: "1px solid",
                                    borderColor: "divider",
                                }}
                            >
                                {key}: {value}
                            </Box>
                        ))}
                    </Box>
                )}
            </Box>
        )
    }

    function renderBadge(text: string, colorType: "success" | "warning" | "error" | "info" | "secondary" | "default") {
        let dotColor = "text.disabled"
        let bgColor = "action.hover"
        let textColor = "text.secondary"

        if (colorType === "success") {
            dotColor = "success.main"
            bgColor = "rgba(46, 125, 50, 0.1)"
            textColor = "success.main"
        } else if (colorType === "warning") {
            dotColor = "warning.main"
            bgColor = "rgba(237, 108, 2, 0.1)"
            textColor = "warning.main"
        } else if (colorType === "error") {
            dotColor = "error.main"
            bgColor = "rgba(211, 47, 47, 0.1)"
            textColor = "error.main"
        } else if (colorType === "info") {
            dotColor = "info.main"
            bgColor = "rgba(2, 136, 209, 0.1)"
            textColor = "info.main"
        } else if (colorType === "secondary") {
            dotColor = "secondary.main"
            bgColor = "rgba(156, 39, 176, 0.1)"
            textColor = "secondary.main"
        }

        return (
            <Box sx={{...SX.badge, bgcolor: bgColor, color: textColor}}>
                <Box sx={{...SX.dot, bgcolor: dotColor}}/>
                {text}
            </Box>
        )
    }

    function renderActions() {
        if (!onUpdate) return undefined
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

    async function handleSave() {
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
}
