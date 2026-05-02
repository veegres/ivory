import {useRouterNodeMetrics} from "../../../api/node/hook"
import {NodeMetrics, SshConnection} from "../../../api/node/type"
import {ErrorSmart} from "../../view/box/ErrorSmart"
import {HistoryTrackerChart} from "../../view/chart/HistoryTrackerChart"
import {MonitorLoading} from "./MonitorLoading"
import {MonitorRow} from "./MonitorRow"

type Props = {
    connection: SshConnection,
    interval?: number,
}

export function MonitorSystem(props: Props) {
    const {connection, interval = 1000 * 3} = props
    const metrics = useRouterNodeMetrics(connection, interval)

    if (metrics.isError) return <ErrorSmart error={metrics.error}/>
    if (metrics.isPending) return <MonitorLoading count={4}/>

    return (
        <MonitorRow>
            <HistoryTrackerChart
                label={"CPU Usage"}
                unit={"%"}
                data={metrics.data}
                selector={getCpuUsageDelta}
                color={"#3f51b5"}
                min={0}
                max={100}
            />
            <HistoryTrackerChart
                label={"Memory Usage"}
                unit={"%"}
                data={metrics.data}
                selector={(m) => m.memory.usagePercent}
                color={"#4caf50"}
                min={0}
                max={100}
            />
            <HistoryTrackerChart
                label={"Network Download"}
                unit={"KB/s"}
                data={metrics.data}
                selector={getNetRxDelta}
                color={"#ff9800"}
            />
            <HistoryTrackerChart
                label={"Network Upload"}
                unit={"KB/s"}
                data={metrics.data}
                selector={getNetTxDelta}
                color={"#9c27b0"}
            />
        </MonitorRow>
    )

    function getCpuUsageDelta(l: NodeMetrics, p?: NodeMetrics) {
        if (!p) return undefined
        const totalDiff = l.cpu.totalTicks - p.cpu.totalTicks
        const idleDiff = l.cpu.idleTicks - p.cpu.idleTicks
        return totalDiff > 0 ? (totalDiff - idleDiff) / totalDiff * 100 : 0
    }

    function getNetRxDelta(l: NodeMetrics, p?: NodeMetrics) {
        if (!p) return undefined
        const rx = (l.network.receivedBytes - p.network.receivedBytes) / 1024 / (interval / 1000)
        return rx < 0 ? 0 : rx
    }

    function getNetTxDelta(l: NodeMetrics, p?: NodeMetrics) {
        if (!p) return undefined
        const tx = (l.network.transmittedBytes - p.network.transmittedBytes) / 1024 / (interval / 1000)
        return tx < 0 ? 0 : tx
    }
}
