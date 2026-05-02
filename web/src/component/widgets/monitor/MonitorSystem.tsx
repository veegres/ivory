import {Box} from "@mui/material"
import {useEffect, useState} from "react"

import {useRouterNodeMetrics} from "../../../api/node/hook"
import {NodeMetrics, SshConnection} from "../../../api/node/type"
import {SxPropsMap} from "../../../app/type"
import {ErrorSmart} from "../../view/box/ErrorSmart"
import {LineChartHistory} from "../../view/chart/LineChartHistory"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1, width: "100%"},
    row: {display: "flex", gap: 1, width: "100%"},
    item: {
        flex: 1, display: "flex", flexDirection: "column", borderRadius: 1, padding: "8px 10px 0px 10px",
        border: "1px solid", borderColor: "divider", minWidth: 0,
    },
    head: {display: "flex", justifyContent: "space-between", alignItems: "center", mb: 0.5},
    label: {fontSize: "12px", fontFamily: "monospace"},
    value: {fontSize: "12px", color: "text.secondary"},
}

type Props = {
    connection: SshConnection,
    historyLength?: number,
    interval?: number,
}

interface DisplayMetrics {
    cpu: number,
    mem: number,
    netRx: number,
    netTx: number,
}

export function MonitorSystem(props: Props) {
    const {connection, historyLength = 60, interval = 1000} = props
    const [history, setHistory] = useState<DisplayMetrics[]>([])
    const [lastRaw, setLastRaw] = useState<NodeMetrics | null>(null)

    const metrics = useRouterNodeMetrics(connection, interval)

    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(handleEffectMetrics, [metrics.data])

    if (metrics.error) return <ErrorSmart error={metrics.error}/>

    return (
        <Box sx={SX.box}>
            <Box sx={SX.row}>
                {renderChart("CPU Usage (%)", history.map(h => h.cpu), "#3f51b5", 0, 100)}
                {renderChart("Memory Usage (%)", history.map(h => h.mem), "#4caf50", 0, 100)}
            </Box>
            <Box sx={SX.row}>
                {renderChart("Network RX (KB/s)", history.map(h => h.netRx), "#ff9800")}
                {renderChart("Network TX (KB/s)", history.map(h => h.netTx), "#9c27b0")}
            </Box>
        </Box>
    )

    function renderChart(label: string, data: (number | null)[], color: string, min?: number, max?: number) {
        const paddedData = [...Array(Math.max(0, historyLength - data.length)).fill(null), ...data]
        const latestValue = data.length > 0 ? data[data.length - 1] : undefined

        return (
            <Box sx={SX.item}>
                <Box sx={SX.head}>
                    <Box sx={SX.label}>{label}</Box>
                    {latestValue && (<Box sx={SX.value}>{latestValue.toFixed(2)}</Box>)}
                </Box>
                <LineChartHistory data={paddedData} color={color} min={min} max={max}/>
            </Box>
        )
    }

    function handleEffectMetrics() {
        if (metrics.data) {
            const raw = metrics.data
            let netRx = 0
            let netTx = 0

            if (lastRaw) {
                const durationSeconds = interval / 1000
                netRx = (raw.network.receivedBytes - lastRaw.network.receivedBytes) / 1024 / durationSeconds
                netTx = (raw.network.transmittedBytes - lastRaw.network.transmittedBytes) / 1024 / durationSeconds

                if (netRx < 0) netRx = 0
                if (netTx < 0) netTx = 0
            }

            const display: DisplayMetrics = {
                cpu: raw.cpu.usagePercent,
                mem: raw.memory.usagePercent,
                netRx,
                netTx,
            }

            setHistory(prev => {
                const next = [...prev, display]
                if (next.length > historyLength) return next.slice(next.length - historyLength)
                return next
            })
            setLastRaw(raw)
        }
    }
}
