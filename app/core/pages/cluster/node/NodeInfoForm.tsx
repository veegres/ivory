import {Box, TextField} from "@mui/material"
import {useState} from "react"

import {Node, NodeConfig} from "../../../../features/cluster/type"
import {CancelIconButton, EditIconButton, SaveIconButton} from "../../../../shared/component/button/IconButtons"
import {SxPropsMap} from "../../../../shared/helper/type"

const SX: SxPropsMap = {
    box: {display: "flex", width: "100%", padding: "14px 10px", gap: 2},
    body: {display: "flex", flexDirection: "column", width: "100%", gap: 2},
    container: {
        display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(200px, 1fr))", gap: 2,
        alignItems: "center", flex: 1,
    },
    actions: {display: "flex", flexDirection: "column"},
}

type Props = {
    node: Node,
    onUpdate?: (config: NodeConfig) => void,
}

export function NodeInfoForm(props: Props) {
    const {node, onUpdate} = props

    const [isEditing, setIsEditing] = useState(false)
    const [config, setConfig] = useState<NodeConfig>(node.config)
    const [loading, setLoading] = useState(false)

    return (
        <Box sx={SX.box}>
            <Box sx={SX.body}>
                <Box sx={SX.container}>
                    {renderItem("Host", node.config.host, "host", "text")}
                    {renderItem("SSH", node.config.sshPort, "sshPort", "number")}
                    {renderItem("Keeper", node.config.keeperPort, "keeperPort", "number")}
                    {renderItem("Database", node.config.dbPort, "dbPort", "number")}
                </Box>
                <Box sx={SX.container}>
                    {renderState()}
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
                value={isEditing ? (config[field] ?? "") : (value ?? "")}
                disabled={!isEditing || loading}
                onChange={e => {
                    const val = e.target.value
                    const parsed = type === "number" ? (val === "" ? undefined : parseInt(val)) : val
                    setConfig({...config, [field]: parsed})
                }}
            />
        )
    }

    function renderState() {
        return (
            <TextField
                fullWidth
                size={"small"}
                label={"State"}
                value={node.keeper.state ?? "-"}
                disabled={true}
            />
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
            onUpdate?.(config)
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
