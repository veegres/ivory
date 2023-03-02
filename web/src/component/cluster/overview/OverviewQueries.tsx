import {TabProps} from "./Overview";
import {Query} from "../../shared/query/Query";
import {QueryType, SxPropsMap} from "../../../app/types";
import {Box, ToggleButton, ToggleButtonGroup} from "@mui/material";
import {useState} from "react";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", alignItems: "center", gap: 1},
    query: {width: "100%"},
    group: {padding: "0px 10%"},
}

export function OverviewQueries(props: TabProps){
    const {cluster, defaultInstance} = props.info
    const [tab, setTab] = useState(QueryType.STATISTIC)

    return (
        <Box sx={SX.box}>
            <ToggleButtonGroup sx={SX.group} size={"small"} color={"secondary"} value={tab} fullWidth>
                <ToggleButton value={QueryType.REPLICATION} onClick={() => setTab(QueryType.REPLICATION)}>
                    REPLICATION
                </ToggleButton>
                <ToggleButton value={QueryType.STATISTIC} onClick={() => setTab(QueryType.STATISTIC)}>
                    STATISTIC
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
