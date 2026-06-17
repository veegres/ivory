import {Alert, Box, Collapse, Divider, Tab, Tabs} from "@mui/material"
import {useMemo, useState} from "react"

import {useRouterClusterList, useRouterClusterOverview} from "../../../../features/cluster/api/hook"
import {Feature} from "../../../../features/feature"
import {ManageAccess} from "../../../../features/management/component/ManageAccess"
import {AlertCentered} from "../../../../shared/component/box/AlertCentered"
import {ErrorSmart} from "../../../../shared/component/box/ErrorSmart"
import {PageMainBox} from "../../../../shared/component/box/PageMainBox"
import {SxPropsMap} from "../../../../shared/helper/type"
import {getMainKeeper} from "../../../../shared/helper/utils"
import {useStore} from "../../../../shared/provider/StoreProvider"
import {OverviewAction} from "./OverviewAction"
import {OverviewNodes} from "./OverviewNodes"
import {OverviewOptions} from "./OverviewOptions"

const SX: SxPropsMap = {
    headBox: {display: "flex", justifyContent: "space-between", alignItems: "center"},
    infoBox: {padding: "5px 0"},
    chip: {margin: "auto 0", borderRadius: "4px"},
    settingsBox: {height: "100%", display: "flex", flexDirection: "row"},
    mainBox: {display: "flex", flexDirection: "row"},
    leftMainBlock: {flexGrow: 1, overflowX: "auto"},
    dividerVertical: {margin: "0 10px"},
    dividerHorizontal: {margin: "10px 0"},
    collapse: {height: "100%"},
}

export function Overview() {
    const activeTags = useStore(s => s.activeTags)
    const activeCluster = useStore(s => s.activeCluster)
    const manualKeeper = useStore(s => s.manualKeeper)

    const [infoOpen, setInfoOpen] = useState(false)
    const [settingsOpen, setSettingsOpen] = useState(false)

    const clusters = useRouterClusterList(activeTags, false)
    const overview = useRouterClusterOverview(activeCluster?.name, false)

    const [mainDomain, mainNode] = useMemo(
        () => getMainKeeper(overview.data?.nodes, manualKeeper),
        [overview.data?.nodes, manualKeeper],
    )

    return (
        <ManageAccess feature={Feature.ViewNodeDbOverview}>
            <PageMainBox withPadding visible={!!activeCluster || !!clusters.data?.length}>
                <Box sx={SX.headBox}>
                    <Tabs value={0} role={"tab"}>
                        <Tab value={0} label={"Overview"}/>
                    </Tabs>
                    {renderActions()}
                </Box>
                <Box sx={SX.infoBox}>{renderInfoBlock()}</Box>
                <Box sx={SX.mainBox}>
                    <Box sx={SX.leftMainBlock}>{renderMainBlock()}</Box>
                    <Box>{renderSettingsBlock()}</Box>
                </Box>
            </PageMainBox>
        </ManageAccess>
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
            selectOptions={settingsOpen}
            toggleOptions={() => setSettingsOpen(!settingsOpen)}
        />
    }

    function renderInfoBlock() {
        return (
            <Collapse in={infoOpen}>
                <Alert severity={"info"} onClose={() => setInfoOpen(false)}>{renderInfo()}</Alert>
                <Divider sx={SX.dividerHorizontal} orientation={"horizontal"} flexItem/>
            </Collapse>
        )
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

    function renderSettingsBlock() {
        if (!activeCluster) return
        return (
            <Collapse sx={SX.collapse} in={settingsOpen} orientation={"horizontal"} unmountOnExit>
                <Box sx={SX.settingsBox}>
                    <Divider sx={SX.dividerVertical} orientation={"vertical"} flexItem/>
                    <OverviewOptions
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
