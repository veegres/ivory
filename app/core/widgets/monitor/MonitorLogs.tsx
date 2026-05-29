import {useRouterNodePlatformLogs} from "../../../features/node/hook"
import {PlatformConnection} from "../../../features/node/type"
import {ChartBox} from "../../../shared/component/box/ChartBox"
import {Logs} from "../logs/Logs"

type Props = {
    connection: PlatformConnection,
}

export function MonitorLogs(props: Props) {
    const {connection} = props
    const request = {connection, name: connection.host, tail: 50}
    const logs = useRouterNodePlatformLogs(request)

    return (
        <ChartBox label={"Logs"} value={logs.data.length} width={"100%"} fixed={false} unit={"rows"}>
            <Logs logs={logs.data} loading={logs.isFetching}/>
        </ChartBox>
    )
}