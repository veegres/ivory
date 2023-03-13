import {Database, SxPropsMap} from "../../../type/common";
import {useState} from "react";
import {QueryType} from "../../../type/query";
import {Box, Collapse, ToggleButton, ToggleButtonGroup} from "@mui/material";
import {QueryNew} from "../../shared/query/QueryNew";
import {Query} from "../../shared/query/Query";

const SX: SxPropsMap = {
    box: {display: "flex", gap: 3},
    filters: {display: "flex", flexDirection: "column", alignItems: "center", gap: 2},
    query: {flexGrow: 1, overflow: "hidden"},
    group: {margin: "0px 5px", width: "100%"},
}

type Props = {
    cluster: string,
    db: Database,
    showNew: boolean,
}

export function InstanceMainQueries(props: Props){
    const {cluster, db, showNew} = props
    const [tab, setTab] = useState(QueryType.ACTIVITY)

    return (
        <Box sx={SX.box}>
            <Box sx={SX.filters}>
                <ToggleButtonGroup sx={SX.group} size={"small"} color={"secondary"} value={tab} orientation={"vertical"}>
                    <ToggleButton value={QueryType.ACTIVITY} onClick={() => setTab(QueryType.ACTIVITY)}>
                        ACTIVITY
                    </ToggleButton>
                    <ToggleButton value={QueryType.STATISTIC} onClick={() => setTab(QueryType.STATISTIC)}>
                        STATISTIC
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
                <Collapse in={showNew}><QueryNew type={tab}/></Collapse>
                <Query type={tab} cluster={cluster} db={db}/>
            </Box>
        </Box>
    )
}
