import {useRouterNodePlatformLogs} from "../../../features/node/hook"
import {PlatformConnection} from "../../../features/node/type"
import {ChartBox} from "../../../shared/component/box/ChartBox"
import {Logs} from "../logs/Logs"

type Props = {
    connection: PlatformConnection,
}

export function MonitorLogs(props: Props) {
    const {connection} = props
    const r = {connection, name: connection.host, tail: 50}
    const logs = useRouterNodePlatformLogs(r)
    const l = logs.data ?? []

    return (
        <ChartBox label={"Logs"} value={l.length} width={"100%"} fixed={false} unit={"rows"}>
            <Logs logs={l} loading={logs.isFetching}/>
        </ChartBox>
    )
}