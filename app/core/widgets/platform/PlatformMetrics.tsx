import {Box} from "@mui/material"
import {useEffect, useState} from "react"

import {useRouterNodeMetrics} from "../../../features/node/hook"
import {PlatformConnection, PlatformMetricsResponse as NodeMetrics} from "../../../features/node/type"
import {ErrorSmart} from "../../../shared/component/box/ErrorSmart"
import {HistoryTrackerChart} from "../../../shared/component/chart/HistoryTrackerChart"
import {SkeletonGroup} from "../../../shared/component/progress/SkeletonGroup"
import {SxPropsMap} from "../../../shared/helper/type"

const SX: SxPropsMap = {
    box: {display: "flex", flexWrap: "wrap", justifyContent: "space-between", gap: 2},
}

type Props = {
    connection: PlatformConnection,
    interval?: number,
}

export function PlatformMetrics(props: Props) {
    const {connection, interval = 1000 * 3} = props
    const [cachedError, setCachedError] = useState<Error>()
    const metrics = useRouterNodeMetrics(connection, interval)

    useEffect(() => {
        if (metrics.data) setCachedError(undefined)
        if (metrics.error) setCachedError(metrics.error)
    }, [metrics.error, metrics.data])

    if (cachedError) return <ErrorSmart error={cachedError}/>

    return (
        <Box sx={SX.box}>
            {renderBody()}
        </Box>
    )
    
    function renderBody() {
        if (metrics.isLoading) return <SkeletonGroup count={4}/>

        return (
            <>
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
                    selector={getMemoryUsage}
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
            </>
        )
    }

    function getCpuUsageDelta(l: NodeMetrics, p?: NodeMetrics) {
        if (!p) return undefined
        const totalDiff = l.cpu.totalTicks - p.cpu.totalTicks
        const idleDiff = l.cpu.idleTicks - p.cpu.idleTicks
        return totalDiff > 0 ? (totalDiff - idleDiff) / totalDiff * 100 : 0
    }

    function getMemoryUsage(m: NodeMetrics) {
        const used = m.memory.totalBytes - m.memory.availableBytes
        return used / m.memory.totalBytes * 100
    }

    // We use the configured interval as the elapsed time denominator
    // rather than measuring actual time between samples. This may
    // slightly underestimate throughput under jitter or tab
    // backgrounding, which is acceptable for a low-frequency dashboard.
    function getNetRxDelta(l: NodeMetrics, p?: NodeMetrics) {
        if (!p) return undefined
        const rx = (l.network.receivedBytes - p.network.receivedBytes) / 1024 / (interval / 1000)
        return rx < 0 ? 0 : rx
    }

    // We use the configured interval as the elapsed time denominator
    // rather than measuring actual time between samples. This may
    // slightly underestimate throughput under jitter or tab
    // backgrounding, which is acceptable for a low-frequency dashboard.
    function getNetTxDelta(l: NodeMetrics, p?: NodeMetrics) {
        if (!p) return undefined
        const tx = (l.network.transmittedBytes - p.network.transmittedBytes) / 1024 / (interval / 1000)
        return tx < 0 ? 0 : tx
    }
}
