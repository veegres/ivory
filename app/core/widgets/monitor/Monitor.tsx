import {Box} from "@mui/material"

import {PlatformConnection} from "../../../features/node/type"
import {ChartType, Connection as QueryConnection} from "../../../features/query/type"
import {ErrorSshMissing} from "../../../shared/component/box/ErrorManual"
import {SxPropsMap} from "../../../shared/helper/type"
import {MonitorCommon} from "./MonitorCommon"
import {MonitorContainer} from "./MonitorContainer"
import {MonitorDatabase} from "./MonitorDatabase"
import {MonitorPlatform} from "./MonitorPlatform"
import {MonitorRow} from "./MonitorRow"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", justifyContent: "center", gap: 2},
    info: {display: "flex"},
}

const Charts = {
    common: [ChartType.Databases, ChartType.Connections, ChartType.DatabaseSize, ChartType.DatabaseUptime],
    database: [ChartType.Schemas, ChartType.TablesSize, ChartType.IndexesSize, ChartType.TotalSize],
}

type Props = {
    queryCon?: QueryConnection,
    platformCon?: PlatformConnection,
}

export function Monitor(props: Props) {
    const {queryCon, platformCon} = props

    return (
        <Box sx={SX.box}>
            {renderPlatformInfo()}
            {renderDbInfo()}
        </Box>
    )

    function renderPlatformInfo() {
        if (!platformCon) return <ErrorSshMissing/>
        return (
            <>
                <MonitorPlatform connection={platformCon}/>
                <MonitorContainer connection={platformCon}/>
            </>
        )
    }

    function renderDbInfo() {
        if (!queryCon) return
        return (
            <>
                <MonitorRow>
                    {Charts.common.map(t => (
                        <MonitorCommon key={t} type={t} connection={queryCon}/>
                    ))}
                </MonitorRow>
                <MonitorRow label={queryCon.db.name && `${queryCon.db.name}`}>
                    {Charts.database.map(t => (
                        <MonitorDatabase key={t} type={t} connection={queryCon}/>
                    ))}
                </MonitorRow>
            </>
        )
    }
}
