import {Logs} from "../../../../shared/component/box/Logs"
import {useRouterNodePlatformContainerLogs} from "../../api/NodeHook"
import {PlatformVaultConnection} from "../../api/NodeType"

type Props = {
    connection: PlatformVaultConnection,
    name: string,
}

export function ContainerLogs(props: Props) {
    const {connection, name} = props
    const logs = useRouterNodePlatformContainerLogs({connection, path: name, tail: 50, follow: true})

    return (
        <Logs logs={logs.data} loading={logs.isFetching} reconnect={logs.reconnect}/>
    )
}
