import {ChartItem, Color} from "./ChartItem";
import {QueryChartType, QueryConnection} from "../../../type/query";
import {useRouterQueryChart} from "../../../router/query";

type Props = {
    type: QueryChartType,
    connection: QueryConnection,
}

export function ChartDatabase(props: Props) {
    const {connection: {db}, type} = props

    const chart = useRouterQueryChart(props, !!db.name)

    if (!db.name) return null

    return (
        <ChartItem
            label={chart.data?.name ?? type}
            value={chart.data?.value}
            loading={chart.isPending}
            color={Color.DEEP_PURPLE}
        />
    )
}
