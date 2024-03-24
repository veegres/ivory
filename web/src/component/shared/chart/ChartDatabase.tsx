import {useQuery} from "@tanstack/react-query";
import {QueryApi} from "../../../app/api";
import {ErrorSmart} from "../../view/box/ErrorSmart";
import {ChartItem, Color} from "./ChartItem";
import {Database, SxPropsMap} from "../../../type/common";
import {ChartLoading} from "./ChartLoading";
import {Box} from "@mui/material";
import {QueryChartType} from "../../../type/query";

const SX: SxPropsMap = {
    error: {flexGrow: 1},
}

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
    if (chart.error) return <Box sx={SX.error}><ErrorSmart error={chart.error}/></Box>
    if (chart.isPending) return <ChartLoading count={4}/>

    return (
        <ChartItem label={chart.data.name} value={chart.data.value} color={Color.DEEP_PURPLE}/>
    )
}
