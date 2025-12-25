import {AxiosError} from "axios"

import {ConnectionRequest, QueryChartType} from "../../../api/postgres"
import {useRouterQueryChart} from "../../../api/query/hook"
import {ChartItem, Color} from "./ChartItem"

type Props = {
    type: QueryChartType,
    connection: ConnectionRequest,
}

export function ChartCommon(props: Props) {
    const {type, connection} = props
    const req = {type, connection: {...connection, db: {...connection.db, name: "postgres"}}}
    const chart = useRouterQueryChart(req)

    return (
        <ChartItem
            label={chart.data?.name ?? type}
            value={chart.data?.value}
            loading={chart.isFetching}
            color={Color.INDIGO}
            error={chart.error as AxiosError}
            onClick={() => chart.refetch()}
        />
    )
}
