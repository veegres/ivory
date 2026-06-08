import {Box, TextField} from "@mui/material"
import {useState} from "react"

import {Node, NodeConfig} from "../../../../features/cluster/type"
import {CancelIconButton, EditIconButton, SaveIconButton} from "../../../../shared/component/button/IconButtons"
import {SxPropsMap} from "../../../../shared/helper/type"

const SX: SxPropsMap = {
    table: {
        position: "relative",
        display: "flex", flexWrap: "wrap", gap: 2, width: "100%", border: 1,
        borderRadius: 1, borderColor: "divider", padding: 1, paddingRight: 5,
        alignItems: "center",
    },
    item: {
        display: "flex", gap: 1, alignItems: "center", minHeight: "40px",
    },
    title: {color: "text.disabled", fontWeight: "bold", whiteSpace: "nowrap"},
    value: {fontWeight: "medium", whiteSpace: "nowrap"},
    input: {
        "& .MuiInputBase-root": {height: "32px", fontSize: "0.9rem"},
        "& .MuiInputBase-input": {padding: "4px 8px"},
        width: "120px",
    },
    actions: {position: "absolute", top: 4, right: 4, display: "flex", flexDirection: "column"},
}

type Props = {
    node: Node,
    onUpdate?: (config: NodeConfig) => void,
}

export function NodeInfoTable(props: Props) {
    const {node, onUpdate} = props

    const [isEditing, setIsEditing] = useState(false)
    const [config, setConfig] = useState<NodeConfig>(node.config)
    const [loading, setLoading] = useState(false)

    return (
        <Box sx={SX.table}>
            {renderItem("Host", node.config.host, "host", "text")}
            {renderItem("SSH", node.config.sshPort, "sshPort", "number")}
            {renderItem("Keeper", node.config.keeperPort, "keeperPort", "number")}
            {renderItem("Database", node.config.dbPort, "dbPort", "number")}

            {!isEditing && (
                <Box sx={SX.item}>
                    <Box sx={SX.title}>State</Box>
                    <Box sx={SX.value}>{node.keeper.state}</Box>
                </Box>
            )}

            {renderActions()}
        </Box>
    )

    function renderItem(label: string, value: string | number | undefined, field: keyof NodeConfig, type: "text" | "number") {
        return (
            <Box sx={SX.item}>
                <Box sx={SX.title}>{label}</Box>
                {isEditing ? (
                    <TextField
                        sx={SX.input}
                        size={"small"}
                        type={type}
                        value={config[field] ?? ""}
                        onChange={e => {
                            const val = e.target.value
                            const parsed = type === "number" ? (val === "" ? undefined : parseInt(val)) : val
                            setConfig({...config, [field]: parsed})
                        }}
                        disabled={loading}
                    />
                ) : (
                    <Box sx={SX.value}>{value?.toString() ?? "-"}</Box>
                )}
            </Box>
        )
    }

    function renderActions() {
        if (!onUpdate) return undefined
        return (
            <Box sx={SX.actions}>
                {isEditing ? (
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
        setLoading(true)
        try {
            await onUpdate?.(config)
            setIsEditing(false)
        } finally {
            setLoading(false)
        }
    }

    function handleCancel() {
        setIsEditing(false)
        setConfig(node.config)
    }

    function handleEdit() {
        setIsEditing(true)
        setConfig(node.config)
    }
}
