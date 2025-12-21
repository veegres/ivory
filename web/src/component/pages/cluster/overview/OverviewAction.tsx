import {InfoOutlined, Settings} from "@mui/icons-material"
import {Box, ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material"

import {ActiveCluster, Instance} from "../../../../api/cluster/type"
import {SxPropsMap} from "../../../../app/type"
import {getSidecarConnection} from "../../../../app/utils"
import {OverviewActionInfo} from "./OverviewActionInfo"
import {OverviewActionStatus} from "./OverviewActionStatus"

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", gap: 1},
    toggleButton: {padding: "3px"},
}

type Props = {
    cluster: ActiveCluster,
    mainInstance?: Instance,
    selectInfo: boolean,
    disableInfo: boolean,
    toggleInfo: () => void,
    selectOptions: boolean,
    toggleOptions: () => void,
}

export function OverviewAction(props: Props) {
    const {cluster, toggleOptions, selectOptions, selectInfo, toggleInfo, disableInfo, mainInstance} = props
    const {name} = cluster.cluster
    const sidecar = mainInstance?.sidecar
    return (
        <Box sx={SX.box}>
            {sidecar && sidecar.status && (
                <OverviewActionStatus
                    status={sidecar.status}
                    cluster={name}
                    request={getSidecarConnection(cluster.cluster, sidecar)}
                />
            )}
            <OverviewActionInfo cluster={cluster.cluster} detectBy={cluster.detectBy} mainInstance={mainInstance}/>
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
