import {useQuery} from "@tanstack/react-query";
import {QueryApi} from "../../../app/api";
import {ErrorSmart} from "../../view/box/ErrorSmart";
import {ChartItem, Color} from "./ChartItem";
import {Database, SxPropsMap} from "../../../type/common";
import {Box, Skeleton} from "@mui/material";
import {QueryChartType} from "../../../type/query";

const SX: SxPropsMap = {
    error: {flexGrow: 1},
}

type Props = {
    type: QueryChartType,
    credentialId: string,
    db: Database,
}

export function ChartCommon(props: Props) {
    const {db, credentialId, type} = props

    const chart = useQuery({
        queryKey: ["query", "chart", type, db.host, db.port, credentialId],
        queryFn: () => QueryApi.chart(props),
        retry: false,
    })

    if (chart.error) return <Box sx={SX.error}><ErrorSmart error={chart.error}/></Box>
    if (chart.isPending) return <Skeleton width={200} height={150}/>

    return (
        <ChartItem label={chart.data.name} value={chart.data.value} color={Color.INDIGO}/>
    )
}
