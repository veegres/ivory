import {AxiosError} from "axios"

import {ConnectionRequest, QueryChartType} from "../../../api/postgres"
import {useRouterQueryChart} from "../../../api/query/hook"
import {ChartItem, Color} from "./ChartItem"

type Props = {
    type: QueryChartType,
    connection: ConnectionRequest,
}

export function ChartDatabase(props: Props) {
    const {connection: {db}, type} = props

    const chart = useRouterQueryChart(props)

    if (!db.name) return null

    return (
        <ChartItem
            label={chart.data?.name ?? type}
            value={chart.data?.value}
            loading={chart.isFetching}
            color={Color.DEEP_PURPLE}
            error={chart.error as AxiosError}
            onClick={() => chart.refetch()}
        />
    )
}
