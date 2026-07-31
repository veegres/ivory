import {Box, FormControl, InputLabel, MenuItem, Select} from "@mui/material"

import {Plugins} from "../../../features/cluster/api/ClusterType"
import {KeeperPlugin} from "../../../features/node/api/NodeType"
import {DbPlugin} from "../../../features/query/api/QueryType"
import {InfoColorBox} from "../../../shared/component/box/InfoColorBox"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {DbPluginOptions, KeeperPluginOptions, ReleaseStageOptions} from "../../../shared/helper/HelperUtils"

const SX: SxPropsMap = {
    field: {flex: "1 1 var(--size-field)", minWidth: "var(--size-field)"},
    item: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 1.5, width: "100%"},
}

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
                    {renderKeeperItem(KeeperPlugin.PATRONI_POSTGRES)}
                    {renderKeeperItem(KeeperPlugin.NATIVE_POSTGRES)}
                    {renderKeeperItem(KeeperPlugin.NATIVE_ETCD)}
                    {renderKeeperItem(KeeperPlugin.NATIVE_REDIS)}
                    {renderKeeperItem(KeeperPlugin.NATIVE_CLICKHOUSE)}
                </Select>
            </FormControl>
            <FormControl sx={SX.field} fullWidth size={"small"}>
                <InputLabel id={"database-plugin"}>Database Plugin</InputLabel>
                <Select
                    labelId={"database-plugin"}
                    label={"Database Plugin"}
                    value={plugins.database}
                    disabled={disabled}
                    onChange={(e) => handleDatabaseUpdate(e.target.value as DbPlugin)}
                >
                    <MenuItem value={DbPlugin.POSTGRES}>{DbPluginOptions[DbPlugin.POSTGRES].label}</MenuItem>
                    <MenuItem value={DbPlugin.ETCD}>{DbPluginOptions[DbPlugin.ETCD].label}</MenuItem>
                    <MenuItem value={DbPlugin.REDIS}>{DbPluginOptions[DbPlugin.REDIS].label}</MenuItem>
                    <MenuItem value={DbPlugin.CLICKHOUSE}>{DbPluginOptions[DbPlugin.CLICKHOUSE].label}</MenuItem>
                </Select>
            </FormControl>
        </>
    )

    function renderKeeperItem(plugin: KeeperPlugin) {
        const option = KeeperPluginOptions[plugin]
        const stage = ReleaseStageOptions[option.stage]
        return (
            <MenuItem key={plugin} value={plugin}>
                <Box sx={SX.item}>
                    <Box>{option.label}</Box>
                    <InfoColorBox label={stage.label} title={stage.description} color={stage.color}/>
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
