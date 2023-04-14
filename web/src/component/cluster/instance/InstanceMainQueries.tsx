import {Database, SxPropsMap} from "../../../type/common";
import {useState} from "react";
import {QueryType} from "../../../type/query";
import {Box, ToggleButton, ToggleButtonGroup} from "@mui/material";
import {Query} from "../../shared/query/Query";

const SX: SxPropsMap = {
    box: {display: "flex", gap: 3},
    filters: {display: "flex", flexDirection: "column", alignItems: "center", gap: 2},
    query: {flexGrow: 1, overflow: "hidden"},
    group: {margin: "0px 5px", width: "100%"},
}

type Props = {
    credentialId?: string,
    db: Database,
}

export function InstanceMainQueries(props: Props){
    const {credentialId, db} = props
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
                <Query type={tab} credentialId={credentialId} db={db}/>
            </Box>
        </Box>
    )
}
