import {AxiosError} from "axios"

import {useRouterQueryChart} from "../../../api/query/hook"
import {ChartType, Connection} from "../../../api/query/type"
import {Color, MonitorItem} from "./MonitorItem"

type Props = {
    type: ChartType,
    connection: Connection,
}

export function MonitorDatabase(props: Props) {
    const {connection: {db}, type} = props

    const chart = useRouterQueryChart(props)

    if (!db.name) return null

    return (
        <MonitorItem
            label={chart.data?.name ?? type}
            value={chart.data?.value}
            loading={chart.isFetching}
            color={Color.DEEP_PURPLE}
            error={chart.error as AxiosError}
            onClick={() => chart.refetch()}
        />
    )
}
