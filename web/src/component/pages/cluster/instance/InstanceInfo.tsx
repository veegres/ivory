import {Box, Paper, ToggleButton, ToggleButtonGroup} from "@mui/material";
import {InstanceInfoStatus} from "./InstanceInfoStatus";
import {InstanceInfoTable} from "./InstanceInfoTable";
import {SxPropsMap} from "../../../../type/general";
import {InstanceTabType, InstanceWeb} from "../../../../type/instance";
import {QueryActivity} from "../../../shared/query/QueryActivity";
import {QueryConnection} from "../../../../type/query";
import {SxPropsFormatter} from "../../../../app/utils";

const SX: SxPropsMap = {
    info: {display: "flex", flexDirection: "column", gap: 1, margin: "5px 0"},
    paper: {padding: "5px", bgcolor: SxPropsFormatter.style.bgImageSelected},
}

type Props = {
    instance: InstanceWeb,
    tab: InstanceTabType,
    onTab: (tab: InstanceTabType) => void,
    connection: QueryConnection,
}

export function InstanceInfo(props: Props) {
    const {instance, tab, onTab, connection} = props

    return (
        <Box sx={SX.info}>
            <ToggleButtonGroup size={"small"} color={"secondary"} fullWidth value={tab}>
                <ToggleButton value={InstanceTabType.CHART} onClick={() => onTab(InstanceTabType.CHART)}>
                    Charts
                </ToggleButton>
                <ToggleButton value={InstanceTabType.QUERY} onClick={() => onTab(InstanceTabType.QUERY)}>
                    Queries
                </ToggleButton>
            </ToggleButtonGroup>
            <InstanceInfoStatus role={instance.role}/>
            <Paper sx={SX.paper} variant={"outlined"}>
                <InstanceInfoTable instance={instance}/>
            </Paper>
            <Paper sx={SX.paper} variant={"outlined"}>
                <QueryActivity connection={connection}/>
            </Paper>
        </Box>
    )
}
