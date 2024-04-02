import {Box, ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material";
import {OverviewActionInfo} from "./OverviewActionInfo";
import {InfoOutlined, Settings} from "@mui/icons-material";
import {ActiveCluster} from "../../../../type/cluster";
import {SxPropsMap} from "../../../../type/general";
import {OverviewActionStatus} from "./OverviewActionStatus";
import {getSidecarConnection} from "../../../../app/utils";

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", gap: 1},
    toggleButton: {padding: "3px"},
}

type Props = {
    cluster: ActiveCluster,
    selectInfo: boolean,
    disableInfo: boolean,
    toggleInfo: () => void,
    selectOptions: boolean,
    toggleOptions: () => void,
}

export function OverviewAction(props: Props) {
    const {cluster, toggleOptions, selectOptions, selectInfo, toggleInfo, disableInfo} = props
    const {name} = cluster.cluster
    const {sidecar} = cluster.defaultInstance
    const request = getSidecarConnection(cluster.cluster, sidecar)
    return (
        <Box sx={SX.box}>
            {sidecar.status && <OverviewActionStatus status={sidecar.status} cluster={name} request={request}/>}
            <OverviewActionInfo cluster={cluster}/>
            <ToggleButtonGroup size={"small"}>
                <ToggleButton
                    sx={SX.toggleButton}
                    value={"settings"}
                    selected={selectOptions}
                    onClick={toggleOptions}
                >
                    <Tooltip title={"Cluster Options"} placement={"top"}><Settings/></Tooltip>
                </ToggleButton>
                <ToggleButton
                    sx={SX.toggleButton}
                    value={"info"}
                    selected={selectInfo}
                    disabled={disableInfo}
                    onClick={toggleInfo}
                >
                    <Tooltip title={"Tab Description"} placement={"top"}><InfoOutlined/></Tooltip>
                </ToggleButton>
            </ToggleButtonGroup>
        </Box>
    )
}
