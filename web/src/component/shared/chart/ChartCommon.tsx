import {useQuery} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {ErrorSmart} from "../../view/box/ErrorSmart";
import {ChartItem, Color} from "./ChartItem";
import {Database, SxPropsMap} from "../../../type/common";
import {ChartLoading} from "./ChartLoading";
import {Box} from "@mui/material";

const SX: SxPropsMap = {
    error: {flexGrow: 1},
}

type Props = {
    credentialId?: string,
    db: Database,
}

export function ChartCommon(props: Props) {
    const {db} = props

    const common = useQuery(
        ["query", "chart", "common", db.host, db.port],
        () => queryApi.chartCommon(props),
        {retry: false})

    if (common.error) return <Box sx={SX.error}><ErrorSmart error={common.error}/></Box>
    if (common.isLoading) return <ChartLoading count={3}/>

    return (
        <>
            {common.data?.map((chart, index) => (
                <ChartItem key={index} label={chart.name} value={chart.value} color={Color.INDIGO}/>
            ))}
        </>
    )
}
