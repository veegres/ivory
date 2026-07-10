import {Box, FormControl, InputLabel, MenuItem, Select} from "@mui/material"

import {Plugins} from "../../../features/cluster/api/ClusterType"
import {KeeperPlugin} from "../../../features/node/api/NodeType"
import {DbPlugin} from "../../../features/query/api/QueryType"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {DbPluginOptions, KeeperPluginOptions} from "../../../shared/helper/HelperUtils"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1.5},
}

type Props = {
    plugins: Plugins,
    onUpdate: (plugins: Plugins) => void,
    disabled?: boolean,
}

export function OptionsPlugins(props: Props) {
    const {plugins, onUpdate, disabled = false} = props

    return (
        <Box sx={SX.box}>
            <FormControl fullWidth size={"small"}>
                <InputLabel id={"keeper-plugin"}>Keeper Plugin</InputLabel>
                <Select
                    labelId={"keeper-plugin"}
                    label={"Keeper Plugin"}
                    value={plugins.keeper}
                    disabled={disabled}
                    onChange={(e) => handleKeeperUpdate(e.target.value as KeeperPlugin)}
                >
                    <MenuItem value={KeeperPlugin.PATRONI_POSTGRES}>{KeeperPluginOptions[KeeperPlugin.PATRONI_POSTGRES].label}</MenuItem>
                    <MenuItem value={KeeperPlugin.NATIVE_POSTGRES}>{KeeperPluginOptions[KeeperPlugin.NATIVE_POSTGRES].label}</MenuItem>
                    <MenuItem value={KeeperPlugin.NATIVE_ETCD}>{KeeperPluginOptions[KeeperPlugin.NATIVE_ETCD].label}</MenuItem>
                </Select>
            </FormControl>
            <FormControl fullWidth size={"small"}>
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
                </Select>
            </FormControl>
        </Box>
    )

    function handleKeeperUpdate(keeper: KeeperPlugin) {
        onUpdate({...plugins, keeper})
    }

    function handleDatabaseUpdate(database: DbPlugin) {
        onUpdate({...plugins, database})
    }
}
