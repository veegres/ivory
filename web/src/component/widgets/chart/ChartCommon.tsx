import {ChartItem, Color} from "./ChartItem";
import {QueryChartType, QueryConnection} from "../../../type/query";
import {useRouterQueryChart} from "../../../router/query";
import {AxiosError} from "axios";

type Props = {
    type: QueryChartType,
    connection: QueryConnection,
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
