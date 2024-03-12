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
    credentialId: string,
    db: Database,
}

export function ChartDatabase(props: Props) {
    const {db, credentialId} = props

    const database = useQuery({
        queryKey: ["query", "chart", "database", db.host, db.port, db.name, credentialId],
        queryFn: () => queryApi.chartDatabase(props),
        retry: false, enabled: !!db.name,
    })

    if (!db.name) return null
    if (database.error) return <Box sx={SX.error}><ErrorSmart error={database.error}/></Box>
    if (database.isPending) return <ChartLoading count={4}/>

    return (
        <>
            {database.data?.map((chart, index) => (
                <ChartItem key={index} label={chart.name} value={chart.value} color={Color.DEEP_PURPLE}/>
            ))}
        </>
    )
}
