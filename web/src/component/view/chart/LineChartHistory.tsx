import {Box} from "@mui/material"
import {LineChart} from "@mui/x-charts/LineChart"

import {SxPropsMap} from "../../../app/type"

const SX: SxPropsMap = {
    box: {width: "100%", height: "100%"},
}

type Props = {
    data: (number | null)[],
    color?: string,
    min?: number,
    max?: number,
    label?: string,
}

export function LineChartHistory(props: Props) {
    const {data, color, min, max} = props

    return (
        <Box sx={SX.box}>
            <LineChart
                series={[{
                    data,
                    color,
                    area: true,
                    showMark: false,
                    valueFormatter: (v) => v !== null ? `${v.toFixed(2)}` : "-",
                }]}
                xAxis={[{
                    data: data.map((_, i) => i),
                    disableLine: true,
                    disableTicks: true,
                    hideTooltip: true,
                    valueFormatter: (v: number) => `${data.length - v}s`,
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
    )
}
