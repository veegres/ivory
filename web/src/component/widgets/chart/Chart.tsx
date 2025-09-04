import {SxPropsMap} from "../../../type/general";
import {Box} from "@mui/material";
import {ChartRow} from "./ChartRow";
import {ChartCommon} from "./ChartCommon";
import {ChartDatabase} from "./ChartDatabase";
import {QueryChartType, QueryConnection} from "../../../type/query";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", justifyContent: "center", gap: 2},
    info: {display: "flex"},
}

const Charts = {
    common: [QueryChartType.Databases, QueryChartType.Connections, QueryChartType.DatabaseSize, QueryChartType.DatabaseUptime],
    database: [QueryChartType.Schemas, QueryChartType.TablesSize, QueryChartType.IndexesSize, QueryChartType.TotalSize],
}

type Props = {
    connection: QueryConnection,
}

export function Chart(props: Props) {
    const {connection: {db}} = props

    return (
        <Box sx={SX.box}>
            <ChartRow>
                {Charts.common.map(t => (
                    <ChartCommon key={t} type={t} connection={props.connection}/>
                ))}
            </ChartRow>
            <ChartRow label={db.name && `${db.name}`}>
                {Charts.database.map(t => (
                    <ChartDatabase key={t} type={t} connection={props.connection}/>
                ))}
            </ChartRow>
        </Box>
    )
}
