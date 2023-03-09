import {TabProps} from "./Overview";
import {Query} from "../../shared/query/Query";
import {QueryType, SxPropsMap} from "../../../app/types";
import {Box, Collapse, Divider, ToggleButton, ToggleButtonGroup} from "@mui/material";
import {useState} from "react";
import {ClusterNoPostgresPassword} from "./OverviewError";
import {QueryNew} from "../../shared/query/QueryNew";

const SX: SxPropsMap = {
    box: {display: "flex", gap: 3},
    filters: {display: "flex", flexDirection: "column", alignItems: "center", gap: 2},
    query: {flexGrow: 1, overflow: "hidden"},
    group: {margin: "0px 5px", width: "100%"},
}

export function OverviewQueries(props: TabProps){
    const {cluster, defaultInstance} = props.info
    const [add, setAdd] = useState(false)
    const [tab, setTab] = useState(QueryType.STATISTIC)

    if (!cluster.credentials.postgresId) return <ClusterNoPostgresPassword/>

    return (
        <Box sx={SX.box}>
            <Box sx={SX.filters}>
                <ToggleButton sx={SX.group} size={"small"} color={"secondary"} value={""} selected={add} onClick={() => setAdd(!add)}>
                    New
                </ToggleButton>
                <Divider flexItem/>
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
            </Box>
            <Box sx={SX.query}>
                <Collapse in={add}><QueryNew type={tab}/></Collapse>
                <Query type={tab} cluster={cluster.name} db={defaultInstance.database}/>
            </Box>
        </Box>
    )
}
