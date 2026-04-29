import {AxiosError} from "axios"

import {useRouterQueryChart} from "../../../api/query/hook"
import {ChartType, Connection} from "../../../api/query/type"
import {Color, MonitorItem} from "./MonitorItem"

type Props = {
    type: ChartType,
    connection: Connection,
}

export function MonitorCommon(props: Props) {
    const {type, connection} = props
    const req = {type, connection: {...connection, db: {...connection.db, name: "postgres"}}}
    const chart = useRouterQueryChart(req)

    return (
        <MonitorItem
            label={chart.data?.name ?? type}
            value={chart.data?.value}
            loading={chart.isFetching}
            color={Color.INDIGO}
            error={chart.error as AxiosError}
            onClick={() => chart.refetch()}
        />
    )
}
