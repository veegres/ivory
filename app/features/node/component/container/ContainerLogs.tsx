import {Box} from "@mui/material"

import {Logs} from "../../../../shared/component/box/Logs"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {useRouterNodePlatformContainerLogs} from "../../api/NodeHook"
import {PlatformVaultConnection} from "../../api/NodeType"

const SX: SxPropsMap = {
    box: {
        display: "flex", flexDirection: "column", gap: 0.5, padding: "5px",
        border: 1, borderRadius: 1, borderColor: "divider",
    },
}

type Props = {
    connection: PlatformVaultConnection,
    name: string,
}

export function ContainerLogs(props: Props) {
    const {connection, name} = props
    const logs = useRouterNodePlatformContainerLogs({connection, path: name, tail: 50, follow: true})

    return (
        <Box sx={SX.box}>
            <Logs logs={logs.data} loading={logs.isFetching} reconnect={logs.reconnect}/>
        </Box>
    )
}
