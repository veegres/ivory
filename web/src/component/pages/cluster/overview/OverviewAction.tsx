import {InfoOutlined, Settings} from "@mui/icons-material"
import {Box, ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material"

import {ActiveCluster, Node} from "../../../../api/cluster/type"
import {SxPropsMap} from "../../../../app/type"
import {getKeeperRequest} from "../../../../app/utils"
import {OverviewActionInfo} from "./OverviewActionInfo"
import {OverviewActionStatus} from "./OverviewActionStatus"

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", gap: 1},
    toggleButton: {padding: "3px"},
}

type Props = {
    cluster: ActiveCluster,
    mainNode: [string?, Node?],
    selectInfo: boolean,
    disableInfo: boolean,
    toggleInfo: () => void,
    selectOptions: boolean,
    toggleOptions: () => void,
}

export function OverviewAction(props: Props) {
    const {cluster, toggleOptions, selectOptions, selectInfo, toggleInfo, disableInfo, mainNode: m} = props
    const {name} = cluster.cluster
    const [mainDomain, mainNode] = m
    const config = mainNode?.config
    const status = mainNode?.keeper.status
    const keeper = config && getKeeperRequest(cluster.cluster, config.host, config.keeperPort)
    return (
        <Box sx={SX.box}>
            {keeper && status && (
                <OverviewActionStatus status={status} cluster={name} request={keeper}/>
            )}
            <OverviewActionInfo cluster={cluster.cluster} manualKeeper={mainDomain}/>
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
