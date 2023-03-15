import {useQuery} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {ErrorAlert} from "../../view/ErrorAlert";
import {ChartItem, Color} from "./ChartItem";
import {Database, SxPropsMap} from "../../../type/common";
import {ChartLoading} from "./ChartLoading";
import {Box} from "@mui/material";

const SX: SxPropsMap = {
    error: {flexGrow: 1},
}

type Props = {
    cluster: string,
    db: Database,
}

export function ChartCommon(props: Props) {
    const {cluster: clusterName, db} = props

    const common = useQuery(
        ["query", "chart", "common", clusterName, db.host, db.port],
        () => queryApi.chartCommon({clusterName, db: {host: db.host, port: db.port}}),
        {retry: false})

    if (common.error) return <Box sx={SX.error}><ErrorAlert error={common.error}/></Box>
    if (common.isLoading) return <ChartLoading count={3}/>

    return (
        <>
            {common.data?.map((chart, index) => (
                <ChartItem
                    key={index}
                    label={chart.name}
                    value={chart.value}
                    color={Color.RED}
                    width={"250px"}
                />
            ))}
        </>
    )
}
