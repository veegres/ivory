import {Box, ToggleButton, ToggleButtonGroup} from "@mui/material"

import {Feature} from "../../../../features/feature"
import {ChartType, Connection as QueryConnection, Type as QueryType} from "../../../../features/query/type"
import {ErrorDbMissing} from "../../../../shared/component/box/ErrorManual"
import {SxPropsMap} from "../../../../shared/helper/type"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"
import {Access} from "../../../widgets/access/Access"
import {MonitorCommon} from "../../../widgets/monitor/MonitorCommon"
import {MonitorDatabase} from "../../../widgets/monitor/MonitorDatabase"
import {MonitorRow} from "../../../widgets/monitor/MonitorRow"
import {Query} from "../../../widgets/query/Query"
import {QueryActivity} from "../../../widgets/query/QueryActivity"
import {QueryConsole} from "../../../widgets/query/QueryConsole"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 2},
    main: {display: "flex", gap: 3},
    filters: {display: "flex", flexDirection: "column", alignItems: "center", gap: 2},
    query: {flexGrow: 1, overflow: "hidden"},
    group: {margin: "0px 5px", width: "100%"},
}

const Charts = {
    common: [ChartType.Databases, ChartType.Connections, ChartType.DatabaseSize, ChartType.DatabaseUptime],
    database: [ChartType.Schemas, ChartType.TablesSize, ChartType.IndexesSize, ChartType.TotalSize],
}

type Props = {
    connection?: QueryConnection,
}

export function NodeMainQueries(props: Props){
    const {connection} = props
    const {queryTab} = useStore(s => s.nodeState)
    const {setQueryTab} = useStoreAction

    if (!connection) return <ErrorDbMissing/>

    return (
        <Box sx={SX.box}>
            <Access feature={Feature.ViewQueryDbInfo}>
                <QueryActivity connection={connection}/>
            </Access>
            <Box sx={SX.main}>
                <Box sx={SX.filters}>
                    <Access feature={Feature.ManageQueryDbConsole}>
                        <ToggleButton
                            sx={SX.group}
                            size={"small"}
                            color={"secondary"}
                            value={QueryType.CONSOLE}
                            selected={queryTab === QueryType.CONSOLE}
                            onClick={() => setQueryTab(QueryType.CONSOLE)}
                        >
                            CONSOLE
                        </ToggleButton>
                    </Access>

                    <Access feature={Feature.ViewQueryDbChart}>
                        <ToggleButton
                            sx={SX.group}
                            size={"small"}
                            color={"secondary"}
                            value={QueryType.CHARTS}
                            selected={queryTab === QueryType.CHARTS}
                            onClick={() => setQueryTab(QueryType.CHARTS)}
                        >
                            CHARTS
                        </ToggleButton>
                    </Access>

                    <ToggleButtonGroup sx={SX.group} size={"small"} color={"secondary"} value={queryTab} orientation={"vertical"}>
                        <ToggleButton value={QueryType.ACTIVITY} onClick={() => setQueryTab(QueryType.ACTIVITY)}>
                            ACTIVITY
                        </ToggleButton>
                        <ToggleButton value={QueryType.STATISTIC} onClick={() => setQueryTab(QueryType.STATISTIC)}>
                            STATISTIC
                        </ToggleButton>
                        <ToggleButton value={QueryType.BLOAT} onClick={() => setQueryTab(QueryType.BLOAT)}>
                            BLOAT
                        </ToggleButton>
                        <ToggleButton value={QueryType.REPLICATION} onClick={() => setQueryTab(QueryType.REPLICATION)}>
                            REPLICATION
                        </ToggleButton>
                        <ToggleButton value={QueryType.OTHER} onClick={() => setQueryTab(QueryType.OTHER)}>
                            OTHER
                        </ToggleButton>
                    </ToggleButtonGroup>
                </Box>
                <Box sx={SX.query}>
                    {queryTab === QueryType.CONSOLE ? (
                        <QueryConsole connection={connection}/>
                    ) : queryTab === QueryType.CHARTS ? (
                        <Box>
                            <MonitorRow>
                                {Charts.common.map(t => (
                                    <MonitorCommon key={t} type={t} connection={connection}/>
                                ))}
                            </MonitorRow>
                            <MonitorRow label={connection.db.name && `${connection.db.name}`}>
                                {Charts.database.map(t => (
                                    <MonitorDatabase key={t} type={t} connection={connection}/>
                                ))}
                            </MonitorRow>
                        </Box>
                    ) : (
                        <Query type={queryTab} connection={connection}/>
                    )}
                </Box>
            </Box>
        </Box>
    )
}
