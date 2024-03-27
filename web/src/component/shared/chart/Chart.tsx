import {Database, SxPropsMap} from "../../../type/general";
import {Box} from "@mui/material";
import {ChartRow} from "./ChartRow";
import {ChartCommon} from "./ChartCommon";
import {ChartDatabase} from "./ChartDatabase";
import {QueryChartType} from "../../../type/query";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", justifyContent: "center", gap: 2},
    info: {display: "flex"},
}

const Charts = {
    common: [QueryChartType.Databases, QueryChartType.Connections, QueryChartType.DatabaseSize, QueryChartType.DatabaseUptime],
    database: [QueryChartType.Schemas, QueryChartType.TablesSize, QueryChartType.IndexesSize, QueryChartType.TotalSize],
}

type Props = {
    credentialId: string,
    db: Database,
}

export function Chart(props: Props) {
    const {credentialId, db} = props

    return (
        <Box sx={SX.box}>
            <ChartRow>
                {Charts.common.map(t => (
                    <ChartCommon key={t} type={t} credentialId={credentialId} db={db}/>
                ))}
            </ChartRow>
            <ChartRow label={db.name && `${db.name}`}>
                {Charts.database.map(t => (
                    <ChartDatabase key={t} type={t} credentialId={credentialId} db={db}/>
                ))}
            </ChartRow>
        </Box>
    )
}
