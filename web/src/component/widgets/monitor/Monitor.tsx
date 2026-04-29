import {Box} from "@mui/material"

import {QueryChartType} from "../../../api/database/type"
import {Connection as NodeConnection} from "../../../api/node/type"
import {Connection as QueryConnection} from "../../../api/query/type"
import {SxPropsMap} from "../../../app/type"
import {MonitorCommon} from "./MonitorCommon"
import {MonitorDatabase} from "./MonitorDatabase"
import {MonitorRow} from "./MonitorRow"
import {MonitorSystem} from "./MonitorSystem"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", justifyContent: "center", gap: 2},
    info: {display: "flex"},
}

const Charts = {
    common: [QueryChartType.Databases, QueryChartType.Connections, QueryChartType.DatabaseSize, QueryChartType.DatabaseUptime],
    database: [QueryChartType.Schemas, QueryChartType.TablesSize, QueryChartType.IndexesSize, QueryChartType.TotalSize],
}

type Props = {
    queryCon: QueryConnection,
    sshCon: NodeConnection,
}

export function Monitor(props: Props) {
    const {queryCon: {db}, sshCon} = props

    return (
        <Box sx={SX.box}>
            <MonitorRow>
                <MonitorSystem connection={sshCon}/>
            </MonitorRow>
            <MonitorRow>
                {Charts.common.map(t => (
                    <MonitorCommon key={t} type={t} connection={props.queryCon}/>
                ))}
            </MonitorRow>
            <MonitorRow label={db.name && `${db.name}`}>
                {Charts.database.map(t => (
                    <MonitorDatabase key={t} type={t} connection={props.queryCon}/>
                ))}
            </MonitorRow>
        </Box>
    )
}
