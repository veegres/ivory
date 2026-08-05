import {Box, FormControl, InputLabel, MenuItem, Select, Tooltip} from "@mui/material"

import {Plugins} from "../../../features/cluster/api/ClusterType"
import {KeeperPlugin} from "../../../features/node/api/NodeType"
import {DbPlugin} from "../../../features/query/api/QueryType"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {DbModelOptions, DbPluginOptions, KeeperPluginOptions, ReleaseStageOptions} from "../../../shared/helper/HelperUtils"

// Mirrors ListKeepers.tsx's plain, muted-text badge style (no colored boxes)
// so the two keeper pickers read as one visual language. Each badge sits in
// a fixed-width column so it lines up across rows regardless of label length.
const SX: SxPropsMap = {
    field: {flex: "1 1 var(--size-field)", minWidth: "var(--size-field)"},
    item: {display: "flex", alignItems: "center", gap: 1.5, width: "100%"},
    label: {flex: "1 1 auto", minWidth: 0, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap"},
    model: {flexShrink: 0, width: "76px", textAlign: "right", fontSize: "12px", color: "primary.main", textTransform: "uppercase"},
    stage: {flexShrink: 0, width: "48px", textAlign: "right", fontSize: "12px", color: "secondary.main", textTransform: "uppercase"},
}

// KEEPER_PLUGIN_ORDER/DB_PLUGIN_ORDER list every plugin grouped by DbModel
// (see the model column each item renders), so the dropdown reads as
// sorted-by-data-model without re-sorting on every render.
const KEEPER_PLUGIN_ORDER: KeeperPlugin[] = [
    KeeperPlugin.PATRONI_POSTGRES,
    KeeperPlugin.NATIVE_POSTGRES,
    KeeperPlugin.NATIVE_MONGO,
    KeeperPlugin.NATIVE_REDIS,
    KeeperPlugin.NATIVE_CLICKHOUSE,
    KeeperPlugin.NATIVE_ETCD,
    KeeperPlugin.NATIVE_ZOOKEEPER,
]

const DB_PLUGIN_ORDER: DbPlugin[] = [
    DbPlugin.POSTGRES,
    DbPlugin.MONGO,
    DbPlugin.REDIS,
    DbPlugin.CLICKHOUSE,
    DbPlugin.ETCD,
    DbPlugin.ZOOKEEPER,
]

type Props = {
    plugins: Plugins,
    onUpdate: (plugins: Plugins) => void,
    disabled?: boolean,
}

export function OptionsPlugins(props: Props) {
    const {plugins, onUpdate, disabled = false} = props

    return (
        <>
            <FormControl sx={SX.field} fullWidth size={"small"}>
                <InputLabel id={"keeper-plugin"}>Keeper Plugin</InputLabel>
                <Select
                    labelId={"keeper-plugin"}
                    label={"Keeper Plugin"}
                    value={plugins.keeper}
                    disabled={disabled}
                    renderValue={(v) => KeeperPluginOptions[v as KeeperPlugin].label}
                    onChange={(e) => handleKeeperUpdate(e.target.value as KeeperPlugin)}
                >
                    {KEEPER_PLUGIN_ORDER.map(renderKeeperItem)}
                </Select>
            </FormControl>
            <FormControl sx={SX.field} fullWidth size={"small"}>
                <InputLabel id={"database-plugin"}>Database Plugin</InputLabel>
                <Select
                    labelId={"database-plugin"}
                    label={"Database Plugin"}
                    value={plugins.database}
                    disabled={disabled}
                    renderValue={(v) => DbPluginOptions[v as DbPlugin].label}
                    onChange={(e) => handleDatabaseUpdate(e.target.value as DbPlugin)}
                >
                    {DB_PLUGIN_ORDER.map(renderDbItem)}
                </Select>
            </FormControl>
        </>
    )

    function renderKeeperItem(plugin: KeeperPlugin) {
        const option = KeeperPluginOptions[plugin]
        const stage = ReleaseStageOptions[option.stage]
        const model = DbModelOptions[option.model]
        return (
            <MenuItem key={plugin} value={plugin}>
                <Box sx={SX.item}>
                    <Box sx={SX.label}>{option.label}</Box>
                    <Tooltip title={model.description} placement={"top"}>
                        <Box sx={SX.model}>{model.label}</Box>
                    </Tooltip>
                    <Tooltip title={stage.description} placement={"top"}>
                        <Box sx={SX.stage}>{stage.label}</Box>
                    </Tooltip>
                </Box>
            </MenuItem>
        )
    }

    function renderDbItem(plugin: DbPlugin) {
        const option = DbPluginOptions[plugin]
        const model = DbModelOptions[option.model]
        return (
            <MenuItem key={plugin} value={plugin}>
                <Box sx={SX.item}>
                    <Box sx={SX.label}>{option.label}</Box>
                    <Tooltip title={model.description} placement={"top"}>
                        <Box sx={SX.model}>{model.label}</Box>
                    </Tooltip>
                </Box>
            </MenuItem>
        )
    }

    function handleKeeperUpdate(keeper: KeeperPlugin) {
        onUpdate({...plugins, keeper})
    }

    function handleDatabaseUpdate(database: DbPlugin) {
        onUpdate({...plugins, database})
    }
}
