import {SxPropsMap} from "../../../../api/management/type";
import {QueryConnection, QueryType} from "../../../../api/query/type";
import {Box, ToggleButton, ToggleButtonGroup} from "@mui/material";
import {Query} from "../../.././widgets/query/Query";
import {QueryConsole} from "../../.././widgets/query/QueryConsole";
import {useStore, useStoreAction} from "../../../../provider/StoreProvider";

const SX: SxPropsMap = {
    box: {display: "flex", gap: 3},
    filters: {display: "flex", flexDirection: "column", alignItems: "center", gap: 2},
    query: {flexGrow: 1, overflow: "hidden"},
    group: {margin: "0px 5px", width: "100%"},
}

type Props = {
    connection: QueryConnection,
}

export function InstanceMainQueries(props: Props){
    const {connection} = props
    const {instance: {queryTab}} = useStore()
    const {setQueryTab} = useStoreAction()

    return (
        <Box sx={SX.box}>
            <Box sx={SX.filters}>
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
                ) : (
                    <Query type={queryTab} connection={connection}/>
                )}
            </Box>
        </Box>
    )
}
