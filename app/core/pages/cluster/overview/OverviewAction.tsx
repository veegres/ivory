import {Cached, InfoOutlined, Settings} from "@mui/icons-material"
import {Box, ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material"

import {useRouterClusterOverview} from "../../../../features/cluster/api/ClusterHook"
import {Cluster, Node} from "../../../../features/cluster/api/ClusterType"
import {SimpleButton} from "../../../../shared/component/button/SimpleButton"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {getKeeperOneRequest} from "../../../../shared/helper/HelperUtils"
import {OverviewActionInfo} from "./OverviewActionInfo"
import {OverviewActionStatus} from "./OverviewActionStatus"
import {OverviewNodesClusterFix} from "./OverviewNodesClusterFix"

const SX: SxPropsMap = {
    status: {
        order: {xs: 3, sm: 2}, flexBasis: {xs: "100%", sm: "auto"},
        display: "flex", justifyContent: "space-between", alignItems: "center", gap: 1, minWidth: 0,
    },
    fixed: {order: {xs: 2, sm: 3}, display: "flex", alignItems: "center", gap: 1},
    toggleButton: {padding: "3px"},
}

type Props = {
    cluster: Cluster,
    mainNode: [string?, Node?],
    selectInfo: boolean,
    disableInfo?: boolean,
    toggleInfo: () => void,
    selectConfig: boolean,
    toggleConfig: () => void,
}

export function OverviewAction(props: Props) {
    const {cluster, toggleConfig, selectConfig, selectInfo, toggleInfo, disableInfo = false, mainNode: m} = props
    const [_, mainNode] = m
    const config = mainNode?.config
    const status = mainNode?.keeper.status
    const keeper = config && getKeeperOneRequest(cluster, config.host, config.keeperPort)
    const overview = useRouterClusterOverview(cluster.name, false)
    return (
        <>
            <Box sx={SX.status}>
                {keeper && status ? (
                    <OverviewActionStatus status={status} cluster={cluster.name} request={keeper}/>
                ) : <Box/>}
                <OverviewActionInfo cluster={cluster}/>
            </Box>
            <Box sx={SX.fixed}>
                <OverviewNodesClusterFix name={cluster.name}/>
                {renderRefresh()}
                <ToggleButtonGroup size={"small"}>
                    <ToggleButton
                        sx={SX.toggleButton}
                        value={"config"}
                        selected={selectConfig}
                        onClick={toggleConfig}
                    >
                        <Tooltip title={"Cluster Config"} placement={"top"}><Settings/></Tooltip>
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
        </>
    )

    function renderRefresh() {
        return (
            <Tooltip title={"Refresh"} placement={"top"} arrow={true}>
                <Box component={"span"}>
                    <SimpleButton loading={overview.isFetching} onClick={() => overview.refetch()}>
                        <Cached/>
                    </SimpleButton>
                </Box>
            </Tooltip>
        )
    }
}
