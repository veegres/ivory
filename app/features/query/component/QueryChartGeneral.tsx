import {AxiosError} from "axios"

import {useRouterQueryChart} from "../api/hook"
import {ChartType, Connection} from "../api/type"
import {Color, QueryChartItem} from "./QueryChartItem"

type Props = {
    type: ChartType,
    connection: Connection,
}

export function QueryChartGeneral(props: Props) {
    const {type, connection} = props
    const req = {type, connection: {...connection, db: {...connection.db, name: "postgres"}}}
    const chart = useRouterQueryChart(req)

    return (
        <QueryChartItem
            label={chart.data?.name ?? type}
            value={chart.data?.value}
            loading={chart.isFetching}
            color={Color.INDIGO}
            error={chart.error as AxiosError}
            onClick={() => chart.refetch()}
        />
    )
}
