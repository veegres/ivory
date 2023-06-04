import {Box, ToggleButton, ToggleButtonGroup} from "@mui/material";
import {InstanceInfoStatus} from "./InstanceInfoStatus";
import {InstanceInfoTable} from "./InstanceInfoTable";
import {SxPropsMap} from "../../../type/common";
import {InstanceTabType, InstanceWeb} from "../../../type/instance";

const SX: SxPropsMap = {
    info: {display: "flex", flexDirection: "column", gap: 2, margin: "5px 0", flex: "0 0 18%"},
}

type Props = {
    instance: InstanceWeb,
    tab: InstanceTabType,
    onTab: (tab: InstanceTabType) => void,
}

export function InstanceInfo(props: Props) {
    const {instance, tab, onTab} = props

    return (
        <Box sx={SX.info}>
            <ToggleButtonGroup size={"small"} color={"secondary"} fullWidth value={tab}>
                <ToggleButton value={InstanceTabType.CHART} onClick={() => onTab(InstanceTabType.CHART)}>Charts</ToggleButton>
                <ToggleButton value={InstanceTabType.QUERY} onClick={() => onTab(InstanceTabType.QUERY)}>Queries</ToggleButton>
            </ToggleButtonGroup>
            <InstanceInfoStatus role={instance.role}/>
            <InstanceInfoTable instance={instance}/>
        </Box>
    )
}
