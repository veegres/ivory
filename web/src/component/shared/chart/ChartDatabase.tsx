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

export function ChartDatabase(props: Props) {
    const {cluster: clusterName, db} = props

    const database = useQuery(
        ["query", "chart", "database", clusterName, db.host, db.port, db.database],
        () => queryApi.chartDatabase({clusterName, db}),
        {retry: false, enabled: !!db.database})

    if (!db.database) return null
    if (database.error) return <Box sx={SX.error}><ErrorAlert error={database.error}/></Box>
    if (database.isLoading) return <ChartLoading count={3}/>

    return (
        <>
            {database.data?.map((chart, index) => (
                <ChartItem
                    key={index}
                    label={chart.name}
                    value={chart.value}
                    color={Color.DEEP_PURPLE}
                    width={"200px"}
                />
            ))}
        </>
    )
}
