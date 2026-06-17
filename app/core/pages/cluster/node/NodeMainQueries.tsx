import {Box, ToggleButton, ToggleButtonGroup} from "@mui/material"

import {Feature} from "../../../../features/feature"
import {ManageAccess} from "../../../../features/management/component/ManageAccess"
import {ChartType, Connection as QueryConnection, Type as QueryType} from "../../../../features/query/api/type"
import {Query} from "../../../../features/query/component/Query"
import {QueryActivity} from "../../../../features/query/component/QueryActivity"
import {QueryChartDatabase} from "../../../../features/query/component/QueryChartDatabase"
import {QueryChartGeneral} from "../../../../features/query/component/QueryChartGeneral"
import {QueryConsole} from "../../../../features/query/component/QueryConsole"
import {DividerBox} from "../../../../shared/component/box/DividerBox"
import {ErrorDbMissing} from "../../../../shared/component/box/ErrorManual"
import {SxPropsMap} from "../../../../shared/helper/type"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 2},
    main: {display: "flex", gap: 3},
    filters: {display: "flex", flexDirection: "column", alignItems: "center", gap: 2},
    query: {flexGrow: 1, overflow: "hidden"},
    group: {margin: "0px 5px", width: "100%"},
}

const CHARTS = {
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
            <ManageAccess feature={Feature.ViewQueryDbInfo}>
                <QueryActivity connection={connection}/>
            </ManageAccess>
            <Box sx={SX.main}>
                <Box sx={SX.filters}>
                    <ManageAccess feature={Feature.ManageQueryDbConsole}>
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
                    </ManageAccess>

                    <ManageAccess feature={Feature.ViewQueryDbChart}>
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
                    </ManageAccess>

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
                            <DividerBox>
                                {CHARTS.common.map(t => (
                                    <QueryChartGeneral key={t} type={t} connection={connection}/>
                                ))}
                            </DividerBox>
                            <DividerBox label={connection.db.name && `${connection.db.name}`}>
                                {CHARTS.database.map(t => (
                                    <QueryChartDatabase key={t} type={t} connection={connection}/>
                                ))}
                            </DividerBox>
                        </Box>
                    ) : (
                        <Query type={queryTab} connection={connection}/>
                    )}
                </Box>
            </Box>
        </Box>
    )
}
