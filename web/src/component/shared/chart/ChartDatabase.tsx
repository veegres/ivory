import {useQuery} from "@tanstack/react-query";
import {QueryApi} from "../../../app/api";
import {ChartItem, Color} from "./ChartItem";
import {Database} from "../../../type/common";
import {QueryChartType} from "../../../type/query";

type Props = {
    type: QueryChartType,
    credentialId: string,
    db: Database,
}

export function ChartDatabase(props: Props) {
    const {db, credentialId, type} = props

    const chart = useQuery({
        queryKey: ["query", "chart", type, db.host, db.port, db.name, credentialId],
        queryFn: () => QueryApi.chart(props),
        retry: false, enabled: !!db.name,
    })

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
