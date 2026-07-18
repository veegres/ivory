import {Alert, Box, Collapse, Divider, Tab, Tabs} from "@mui/material"
import {useMemo, useState} from "react"

import {useRouterClusterList, useRouterClusterOverview} from "../../../../features/cluster/api/ClusterHook"
import {AlertCentered} from "../../../../shared/component/box/AlertCentered"
import {ErrorSmart} from "../../../../shared/component/box/ErrorSmart"
import {PageMainBox} from "../../../../shared/component/box/PageMainBox"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {getMainKeeper} from "../../../../shared/helper/HelperUtils"
import {useStore} from "../../../../shared/provider/StoreProvider"
import {OverviewAction} from "./OverviewAction"
import {OverviewClusterConfig} from "./OverviewClusterConfig"
import {OverviewNodes} from "./OverviewNodes"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    headBox: {display: "flex", flexWrap: "wrap", justifyContent: "space-between", alignItems: "center", gap: 1},
    tabs: {flexGrow: 1},
    chip: {margin: "auto 0", borderRadius: "4px"},
    settingsBox: {height: "100%", display: "flex", flexDirection: "column"},
    mainBox: {display: "flex", flexDirection: "column"},
    leftMainBlock: {flexGrow: 1, overflowX: "auto"},
    divider: {margin: "10px 0", fontSize: "15px", color: "text.secondary"},
    collapse: {height: "100%"},
}

export function Overview() {
    const activeCluster = useStore(s => s.activeCluster)
    const manualKeeper = useStore(s => s.manualKeeper)

    const [infoOpen, setInfoOpen] = useState(false)
    const [configOpen, setConfigOpen] = useState(false)

    const clusters = useRouterClusterList(false)
    const overview = useRouterClusterOverview(activeCluster?.name, false)

    const [mainDomain, mainNode] = useMemo(
        () => getMainKeeper(overview.data?.nodes, manualKeeper),
        [overview.data?.nodes, manualKeeper],
    )

    return (
        <PageMainBox withPadding={true} visible={!!activeCluster || !!clusters.data?.length}>
            <Box sx={SX.box}>
                <Box sx={SX.headBox}>
                    <Tabs sx={SX.tabs} value={0} role={"tab"}>
                        <Tab value={0} label={"Overview"}/>
                    </Tabs>
                    {renderActions()}
                </Box>
                <Collapse in={infoOpen}>
                    <Alert severity={"info"} onClose={() => setInfoOpen(false)}>{renderInfo()}</Alert>
                </Collapse>
                <Box sx={SX.mainBox}>
                    <Box sx={SX.leftMainBlock}>{renderMainBlock()}</Box>
                    <Box>{renderConfigBlock()}</Box>
                </Box>
            </Box>
        </PageMainBox>
    )

    function renderMainBlock() {
        if (!activeCluster) return <AlertCentered text={"Please, select a cluster to see the overview! (click on the name)"}/>
        if (!activeCluster) return <AlertCentered text={"Selected cluster in not in the list"} severity={"warning"}/>
        return <>
            {overview.error && <ErrorSmart error={overview.error}/>}
            <OverviewNodes cluster={activeCluster} nodes={overview.data?.nodes}/>
        </>
    }

    function renderActions() {
        if (!activeCluster) return
        return <OverviewAction
            cluster={activeCluster}
            mainNode={[mainDomain, mainNode]}
            selectInfo={infoOpen}
            toggleInfo={() => setInfoOpen(!infoOpen)}
            selectConfig={configOpen}
            toggleConfig={() => setConfigOpen(!configOpen)}
        />
    }

    function renderInfo() {
        return (
            <>
                The Overview tab offers visibility into the current status of your cluster. From here, you can
                utilize essential features to manage your cluster, including switchover, reinit, restart, reload,
                failover, and more. The leader node is automatically detected by sending requests to each node
                until a successful connection is established. You have the flexibility to change the main node
                to which Ivory sends requests by accessing the settings in the top right corner.
            </>
        )
    }

    function renderConfigBlock() {
        if (!activeCluster) return
        return (
            <Collapse sx={SX.collapse} in={configOpen} unmountOnExit={true}>
                <Box sx={SX.settingsBox}>
                    <Divider sx={SX.divider} textAlign={"left"} flexItem={true}>CONFIGURATION</Divider>
                    <OverviewClusterConfig
                        cluster={activeCluster}
                        overview={overview.data}
                        mainKeeper={mainDomain}
                        manualKeeper={manualKeeper}
                    />
                </Box>
            </Collapse>
        )
    }
}
