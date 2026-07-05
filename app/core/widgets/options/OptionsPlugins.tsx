import {Box, FormControl, InputLabel, MenuItem, Select} from "@mui/material"

import {Plugins} from "../../../features/cluster/api/type"
import {KeeperPlugin} from "../../../features/node/api/type"
import {DbPlugin} from "../../../features/query/api/type"
import {SxPropsMap} from "../../../shared/helper/type"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1.5},
}

type Props = {
    plugins: Plugins,
    onUpdate: (plugins: Plugins) => void,
}

export function OptionsPlugins(props: Props) {
    const {plugins, onUpdate} = props

    return (
        <Box sx={SX.box}>
            <FormControl fullWidth size={"small"}>
                <InputLabel id={"keeper-plugin"}>Keeper Plugin</InputLabel>
                <Select
                    labelId={"keeper-plugin"}
                    label={"Keeper Plugin"}
                    value={plugins.keeper}
                    onChange={(e) => handleKeeperUpdate(e.target.value as KeeperPlugin)}
                >
                    <MenuItem value={KeeperPlugin.PATRONI_POSTGRES}>Patroni Postgres</MenuItem>
                    <MenuItem value={KeeperPlugin.NATIVE_POSTGRES}>Native Postgres</MenuItem>
                    <MenuItem value={KeeperPlugin.NATIVE_ETCD}>Native Etcd</MenuItem>
                </Select>
            </FormControl>
            <FormControl fullWidth size={"small"}>
                <InputLabel id={"database-plugin"}>Database Plugin</InputLabel>
                <Select
                    labelId={"database-plugin"}
                    label={"Database Plugin"}
                    value={plugins.database}
                    onChange={(e) => handleDatabaseUpdate(e.target.value as DbPlugin)}
                >
                    <MenuItem value={DbPlugin.POSTGRES}>Postgres</MenuItem>
                    <MenuItem value={DbPlugin.ETCD}>Etcd</MenuItem>
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
