import {Box} from "@mui/material"

import {SshConnection} from "../../../api/node/type"
import {ChartType, Connection as QueryConnection} from "../../../api/query/type"
import {SxPropsMap} from "../../../app/type"
import {ErrorSmart} from "../../view/box/ErrorSmart"
import {MonitorCommon} from "./MonitorCommon"
import {MonitorDatabase} from "./MonitorDatabase"
import {MonitorRow} from "./MonitorRow"
import {MonitorSystem} from "./MonitorSystem"

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
    sshCon?: SshConnection,
}

export function Monitor(props: Props) {
    const {queryCon, sshCon} = props

    return (
        <Box sx={SX.box}>
            {renderSshInfo()}
            {renderDbInfo()}
        </Box>
    )

    function renderSshInfo() {
        if (!sshCon) return <ErrorSmart error={"provide ssh key and port to see the metrics"}/>
        return (
            <MonitorRow>
                <MonitorSystem connection={sshCon}/>
            </MonitorRow>
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
