import {LineChart} from "@mui/x-charts/LineChart"
import {useEffect, useMemo, useRef, useState} from "react"

import {MetricsResponse as NodeMetrics} from "../../../api/node/type"
import {ChartBox} from "../box/ChartBox"

type Props = {
    label: string,
    data: NodeMetrics | undefined,
    selector: (last: NodeMetrics, penultimate?: NodeMetrics) => number | undefined,
    length?: number,
    color?: string,
    unit?: string,
    min?: number,
    max?: number,
}

export function HistoryTrackerChart(props: Props) {
    const {label, data, selector, length = 60, color, unit, min, max} = props
    const [history, setHistory] = useState<(number | undefined)[]>(new Array(length).fill(undefined))
    const lastRawRef = useRef<NodeMetrics | undefined>(undefined)

    const latestValue = useMemo(handleMemoLatestValue, [history])
    const xData = useMemo(handleMemoXLabels, [length])

    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(handleEffectMetrics, [data])

    return (
        <ChartBox label={label} value={latestValue} unit={unit}>
            <LineChart
                series={[{
                    data: history.map(x => x ?? null),
                    color,
                    area: true,
                    showMark: false,
                    valueFormatter: (v) => v !== null ? `${v.toFixed(2)} ${unit ?? ""}` : "-",
                }]}
                xAxis={[{
                    data: xData,
                    tickLabelStyle: {fontSize: 10},
                    disableLine: true,
                    disableTicks: true,
                    hideTooltip: true,
                    tickMinStep: 10,
                    valueFormatter: (v: number) => `${length - v}s`,
                }]}
                yAxis={[{
                    min, max,
                    position: "none",
                    disableLine: true,
                    disableTicks: true,
                }]}
                slotProps={{legend: {hidden: true} as any}}
                margin={0}
                height={120}
                skipAnimation={true}
            />
        </ChartBox>
    )

    function handleEffectMetrics() {
        if (!data) return
        setHistory(prev => [...prev.slice(1), selector(data, lastRawRef.current)])
        lastRawRef.current = data
    }

    function handleMemoLatestValue() {
        return history[history.length - 1]
    }

    function handleMemoXLabels() {
        return [...Array(length).keys()]
    }
}
