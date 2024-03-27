import {ChartItem, Color} from "./ChartItem";
import {Database} from "../../../type/general";
import {QueryChartType} from "../../../type/query";
import {useRouterQueryChart} from "../../../router/query";

type Props = {
    type: QueryChartType,
    credentialId: string,
    db: Database,
}

export function ChartDatabase(props: Props) {
    const {db, type} = props

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
