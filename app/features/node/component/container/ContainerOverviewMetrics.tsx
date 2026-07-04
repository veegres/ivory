import {Box} from "@mui/material"
import {useEffect, useState} from "react"

import {ErrorSmart} from "../../../../shared/component/box/ErrorSmart"
import {HistoryTrackerChart} from "../../../../shared/component/chart/HistoryTrackerChart"
import {SkeletonGroup} from "../../../../shared/component/progress/SkeletonGroup"
import {SxPropsMap} from "../../../../shared/helper/type"
import {useRouterNodePlatformContainerMetrics} from "../../api/hook"
import {Metrics as ContainerMetricsData, PlatformVaultConnection} from "../../api/type"

const SX: SxPropsMap = {
    box: {display: "flex", flexWrap: "wrap", justifyContent: "space-between", gap: 1},
}

type Props = {
    connection: PlatformVaultConnection,
    name: string,
}

export function ContainerOverviewMetrics(props: Props) {
    const {connection, name} = props
    const [cachedError, setCachedError] = useState<Error>()
    const metrics = useRouterNodePlatformContainerMetrics({connection, name})

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
                    selector={getCpuUsage}
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

    // Container CPU is reported by docker as an already-computed percentage
    // (not a monotonic tick counter like host /proc stats), so it is read
    // directly from the latest sample rather than diffed like getCpuUsageDelta.
    function getCpuUsage(m: ContainerMetricsData) {
        return (m.cpu.totalTicks - m.cpu.idleTicks) / m.cpu.totalTicks * 100
    }

    function getMemoryUsage(m: ContainerMetricsData) {
        const used = m.memory.totalBytes - m.memory.availableBytes
        return used / m.memory.totalBytes * 100
    }

    function getNetRxDelta(l: ContainerMetricsData, p?: ContainerMetricsData, elapsedMs?: number) {
        if (!p || !elapsedMs) return undefined
        const rx = (l.network.receivedBytes - p.network.receivedBytes) / 1024 / (elapsedMs / 1000)
        return rx < 0 ? 0 : rx
    }

    function getNetTxDelta(l: ContainerMetricsData, p?: ContainerMetricsData, elapsedMs?: number) {
        if (!p || !elapsedMs) return undefined
        const tx = (l.network.transmittedBytes - p.network.transmittedBytes) / 1024 / (elapsedMs / 1000)
        return tx < 0 ? 0 : tx
    }
}
