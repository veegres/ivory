import {OverviewInstances} from "./OverviewInstances";
import {OverviewConfig} from "./OverviewConfig";
import {Alert, Box, Collapse, Divider, Link, Tab, Tabs} from "@mui/material";
import {useState} from "react";
import {useStore, useStoreAction} from "../../../../provider/StoreProvider";
import {OverviewBloat} from "./OverviewBloat";
import {InfoAlert} from "../../../view/box/InfoAlert";
import {PageMainBox} from "../../../view/box/PageMainBox";
import {OverviewOptions} from "./OverviewOptions";
import {SxPropsMap} from "../../../../type/general";
import {ActiveCluster, ClusterTabs} from "../../../../type/cluster";
import {OverviewAction} from "./OverviewAction";
import {InstanceWeb} from "../../../../type/instance";
import {useRouterClusterList} from "../../../../router/cluster";

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

const TABS: ClusterTabs = {
    0: {
        label: "Overview",
        body: (cluster: ActiveCluster, activeInstance?: InstanceWeb) => <OverviewInstances info={cluster} activeInstance={activeInstance}/>,
        info: <>
            The Overview tab offers visibility into the current status of your cluster. From here, you can
            utilize essential features to manage your cluster, including switchover, reinit, restart, reload,
            failover, and more. The leader node is automatically detected by sending requests to each instance
            until a successful connection is established. You have the flexibility to change the main instance
            to which Ivory sends requests by accessing the settings in the top right corner.
        </>
    },
    1: {
        label: "Config",
        body: (cluster: ActiveCluster) => <OverviewConfig info={cluster}/>,
        info: <>
            You can adjust your PostgreSQL configurations here, and any changes made will be applied to
            all cluster instances. Instead of rewriting the entire configuration, it applies a patch
            update. If you wish to remove a specific setting, simply set it to <b>null</b>. Keep in mind that
            modifying certain parameters may necessitate restarting PostgreSQL. For further details on how
            this process functions, refer to
            the <Link href={"https://patroni.readthedocs.io/en/latest/rest_api.html#config-endpoint"} target={"_blank"}>documentation</Link>.
        </>
    },
    2: {
        label: "Bloat",
        body: (cluster: ActiveCluster) => <OverviewBloat info={cluster}/>,
        info: <>
            Here, you can efficiently decrease the size of bloated tables and indexes without imposing
            heavy locks. This functionality is powered by
            the <Link href={"https://github.com/dataegret/pgcompacttable"} target={"_blank"}>pgcompacttable</Link>
            tool, which is seamlessly integrated with Ivory for streamlined usage. Ivory simplifies visualization
            and centralizes information about jobs and logs within each cluster, ensuring convenient access when
            needed. It's important to note that this tool can only be executed on the master node, and
            in the target database, the contrib module pgstattuple must be installed using the command
            "<b>CREATE EXTENSION IF NOT EXISTS pgstattuple;</b>".
        </>
    },
}

export function Overview() {
    const {activeCluster, activeClusterTab, activeInstance} = useStore()
    const {setClusterTab} = useStoreAction()
    const [infoOpen, setInfoOpen] = useState(false)
    const [settingsOpen, setSettingsOpen] = useState(false)
    const clusters = useRouterClusterList(undefined, true)
    const tab = TABS[activeClusterTab]

    return (
        <PageMainBox withPadding visible={!!activeCluster || Object.keys(clusters.data ?? {}).length !== 0}>
            <Box sx={SX.headBox}>
                <Tabs value={activeClusterTab} onChange={(_, value) => setClusterTab(value)}>
                    {Object.entries(TABS).map(([key, value]) => (<Tab key={key} label={value.label}/>))}
                </Tabs>
                {renderActions()}
            </Box>
            <Box sx={SX.infoBox}>{renderInfoBlock()}</Box>
            <Box sx={SX.mainBox}>
                <Box sx={SX.leftMainBlock}>{renderMainBlock()}</Box>
                <Box>{renderSettingsBlock()}</Box>
            </Box>
        </PageMainBox>
    )

    function renderMainBlock() {
        if (!activeCluster) return <InfoAlert text={"Please, select a cluster to see the overview! (click on the name)"}/>
        if (!tab) return <InfoAlert text={"Coming soon — we're working on it!"}/>
        return tab.body(activeCluster, activeInstance)
    }

    function renderActions() {
        if (!activeCluster) return null
        return <OverviewAction
            cluster={activeCluster}
            selectInfo={infoOpen}
            disableInfo={!tab.info}
            toggleInfo={() => setInfoOpen(!infoOpen)}
            selectOptions={settingsOpen}
            toggleOptions={() => setSettingsOpen(!settingsOpen)}
        />
    }

    function renderInfoBlock() {
        if (!tab?.info) return null
        return (
            <Collapse in={infoOpen}>
                <Alert severity={"info"} onClose={() => setInfoOpen(false)}>{tab.info}</Alert>
                <Divider sx={SX.dividerHorizontal} orientation={"horizontal"} flexItem/>
            </Collapse>
        )
    }

    function renderSettingsBlock() {
        if (!activeCluster) return null
        return (
            <Collapse sx={SX.collapse} in={settingsOpen} orientation={"horizontal"} unmountOnExit>
                <Box sx={SX.settingsBox}>
                    <Divider sx={SX.dividerVertical} orientation={"vertical"} flexItem/>
                    <OverviewOptions info={activeCluster}/>
                </Box>
            </Collapse>
        )
    }
}
