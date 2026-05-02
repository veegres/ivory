import {Box} from "@mui/material"

import {useRouterNodeMetrics} from "../../../api/node/hook"
import {NodeMetrics, SshConnection} from "../../../api/node/type"
import {SxPropsMap} from "../../../app/type"
import {ErrorSmart} from "../../view/box/ErrorSmart"
import {HistoryTrackerChart} from "../../view/chart/HistoryTrackerChart"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1, width: "100%"},
    row: {display: "flex", gap: 1, width: "100%"},
}

type Props = {
    connection: SshConnection,
    historyLength?: number,
    interval?: number,
}

export function MonitorSystem(props: Props) {
    const {connection, historyLength, interval = 1000} = props
    const metrics = useRouterNodeMetrics(connection, interval)

    if (metrics.error) return <ErrorSmart error={metrics.error}/>

    return (
        <Box sx={SX.box}>
            <Box sx={SX.row}>
                <HistoryTrackerChart
                    label={"CPU Usage"}
                    unit={"%"}
                    data={metrics.data}
                    selector={(m) => m.cpu.usagePercent}
                    color={"#3f51b5"}
                    min={0}
                    max={100}
                    maxLength={historyLength}
                />
                <HistoryTrackerChart
                    label={"Memory Usage"}
                    unit={"%"}
                    data={metrics.data}
                    selector={(m) => m.memory.usagePercent}
                    color={"#4caf50"}
                    min={0}
                    max={100}
                    maxLength={historyLength}
                />
            </Box>
            <Box sx={SX.row}>
                <HistoryTrackerChart
                    label={"Network RX"}
                    unit={"KB/s"}
                    data={metrics.data}
                    selector={getNetRx}
                    color={"#ff9800"}
                    maxLength={historyLength}
                />
                <HistoryTrackerChart
                    label={"Network TX"}
                    unit={"KB/s"}
                    data={metrics.data}
                    selector={getNetTx}
                    color={"#9c27b0"}
                    maxLength={historyLength}
                />
            </Box>
        </Box>
    )

    function getNetRx(data: NodeMetrics, last?: NodeMetrics) {
        if (!last) return 0
        const rx = (data.network.receivedBytes - last.network.receivedBytes) / 1024 / (interval / 1000)
        return rx < 0 ? 0 : rx
    }

    function getNetTx(data: NodeMetrics, last?: NodeMetrics) {
        if (!last) return 0
        const tx = (data.network.transmittedBytes - last.network.transmittedBytes) / 1024 / (interval / 1000)
        return tx < 0 ? 0 : tx
    }
}
