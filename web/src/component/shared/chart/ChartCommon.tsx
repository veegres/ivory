import {ChartItem, Color} from "./ChartItem";
import {Database} from "../../../type/general";
import {QueryChartType} from "../../../type/query";
import {useRouterQueryChart} from "../../../router/query";

type Props = {
    type: QueryChartType,
    credentialId: string,
    db: Database,
}

export function ChartCommon(props: Props) {
    const {type, credentialId, db} = props
    const req = {type, credentialId, db: {...db, name: "postgres"}}
    const chart = useRouterQueryChart(req)

    return (
        <ChartItem
            label={chart.data?.name ?? type}
            value={chart.data?.value}
            loading={chart.isPending}
            color={Color.INDIGO}
        />
    )
}
