import {TabProps} from "./Overview";
import {Query} from "../../shared/query/Query";
import {QueryType, SxPropsMap} from "../../../app/types";
import {Box, ToggleButton, ToggleButtonGroup} from "@mui/material";
import {useState} from "react";
import {ClusterNoPostgresPassword} from "./OverviewError";

const SX: SxPropsMap = {
    box: {display: "flex", gap: 3},
    query: {width: "100%"},
    group: {padding: "0px 5px"},
}

export function OverviewQueries(props: TabProps){
    const {cluster, defaultInstance} = props.info
    const [tab, setTab] = useState(QueryType.STATISTIC)

    if (!cluster.credentials.postgresId) return <ClusterNoPostgresPassword/>

    return (
        <Box sx={SX.box}>
            <ToggleButtonGroup sx={SX.group} size={"small"} color={"secondary"} value={tab} orientation={"vertical"}>
                <ToggleButton value={QueryType.STATISTIC} onClick={() => setTab(QueryType.STATISTIC)}>
                    STATISTIC
                </ToggleButton>
                <ToggleButton value={QueryType.ACTIVITY} onClick={() => setTab(QueryType.ACTIVITY)}>
                    ACTIVITY
                </ToggleButton>
                <ToggleButton value={QueryType.BLOAT} onClick={() => setTab(QueryType.BLOAT)}>
                    BLOAT
                </ToggleButton>
                <ToggleButton value={QueryType.REPLICATION} onClick={() => setTab(QueryType.REPLICATION)}>
                    REPLICATION
                </ToggleButton>
                <ToggleButton value={QueryType.OTHER} onClick={() => setTab(QueryType.OTHER)}>
                    OTHER
                </ToggleButton>
            </ToggleButtonGroup>
            <Box sx={SX.query}>
                <Query type={tab} cluster={cluster.name} db={defaultInstance.database}/>
            </Box>
        </Box>
    )
}
