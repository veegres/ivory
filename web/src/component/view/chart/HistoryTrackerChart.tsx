import {Box} from "@mui/material"
import {LineChart} from "@mui/x-charts/LineChart"
import {useEffect, useMemo, useState} from "react"

import {NodeMetrics} from "../../../api/node/type"
import {SxPropsMap} from "../../../app/type"
import {ChartBox} from "../box/ChartBox"

const SX: SxPropsMap = {
    chart: {width: "100%", height: "100%"},
}

type Props = {
    label: string,
    data: NodeMetrics | undefined,
    selector: (data: NodeMetrics, lastData?: NodeMetrics) => number,
    maxLength?: number,
    color?: string,
    unit?: string,
    min?: number,
    max?: number,
}

export function HistoryTrackerChart(props: Props) {
    const {label, data, selector, maxLength = 60, color, unit, min, max} = props
    const [history, setHistory] = useState<(number | undefined)[]>([])
    const [lastRaw, setLastRaw] = useState<NodeMetrics | undefined>()

    const latestValue = useMemo(handleMemoLatestValue, [history])
    const paddedData = useMemo(handleMemoPaddedData, [history, maxLength])

    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(handleEffectMetrics, [data, maxLength])

    return (
        <ChartBox label={label} value={latestValue} unit={unit}>
            <Box sx={SX.chart}>
                <LineChart
                    series={[{
                        data: paddedData as (number | null)[],
                        color,
                        area: true,
                        showMark: false,
                        valueFormatter: (v) => v !== null ? `${v.toFixed(2)} ${unit ?? ""}` : "-",
                    }]}
                    xAxis={[{
                        data: paddedData.map((_, i) => i),
                        disableLine: true,
                        disableTicks: true,
                        hideTooltip: true,
                        valueFormatter: (v: number) => `${paddedData.length - v}s`,
                    }]}
                    yAxis={[{
                        min,
                        max,
                        position: "none",
                        disableLine: true,
                        disableTicks: true,
                    }]}
                    slotProps={{legend: {hidden: true} as any}}
                    margin={0}
                    height={120}
                />
            </Box>
        </ChartBox>
    )

    function handleEffectMetrics() {
        if (!data) return

        const value = selector(data, lastRaw)
        setHistory(prev => {
            const next = [...prev, value]
            return next.length > maxLength ? next.slice(next.length - maxLength) : next
        })
        setLastRaw(data)
    }

    function handleMemoLatestValue() {
        return history.length > 0 ? history[history.length - 1] : undefined
    }

    function handleMemoPaddedData() {
        return [...Array(Math.max(0, maxLength - history.length)).fill(undefined), ...history]
    }
}
