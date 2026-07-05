import {AxiosError} from "axios"

import {useRouterQueryChart} from "../api/QueryHook"
import {ChartType, Connection} from "../api/QueryType"
import {Color, QueryChartItem} from "./QueryChartItem"

type Props = {
    type: ChartType,
    connection: Connection,
}

export function QueryChartDatabase(props: Props) {
    const {connection: {db}, type} = props
    const chart = useRouterQueryChart(props)

    if (!db.name) return
    return (
        <QueryChartItem
            label={chart.data?.name ?? type}
            value={chart.data?.value}
            loading={chart.isFetching}
            color={Color.DEEP_PURPLE}
            error={chart.error as AxiosError}
            onClick={() => chart.refetch()}
        />
    )
}
