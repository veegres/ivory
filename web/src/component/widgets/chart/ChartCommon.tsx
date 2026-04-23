import {AxiosError} from "axios"

import {QueryChartType} from "../../../api/postgres"
import {useRouterQueryChart} from "../../../api/query/hook"
import {Connection} from "../../../api/query/type"
import {ChartItem, Color} from "./ChartItem"

type Props = {
    type: QueryChartType,
    connection: Connection,
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
