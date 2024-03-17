import {OverviewInstances} from "./OverviewInstances";
import {OverviewConfig} from "./OverviewConfig";
import {Alert, Box, Collapse, Divider, Link, Tab, Tabs} from "@mui/material";
import {useState} from "react";
import {useStore, useStoreAction} from "../../../provider/StoreProvider";
import {OverviewBloat} from "./OverviewBloat";
import {InfoAlert} from "../../view/box/InfoAlert";
import {PageMainBox} from "../../view/box/PageMainBox";
import {useQuery} from "@tanstack/react-query";
import {OverviewOptions} from "./OverviewOptions";
import {SxPropsMap} from "../../../type/common";
import {ActiveCluster, ClusterMap, ClusterTabs} from "../../../type/cluster";
import {OverviewAction} from "./OverviewAction";

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
        body: (cluster: ActiveCluster) => <OverviewInstances info={cluster}/>,
    },
    1: {
        label: "Config",
        body: (cluster: ActiveCluster) => <OverviewConfig info={cluster}/>,
        info: <>
            You can change your postgres configurations here (it will be applied on all cluster nodes).
            It doesn't rewrite all your config, it call patches the existing configuration. If you want to
            remove (reset) some setting just patch it with <i>null</i>. Be aware that some of the parameters
            requires a restart of postgres. More information how it works you can find in
            patroni <Link href={"https://patroni.readthedocs.io/en/latest/SETTINGS.html"} target={"_blank"}>site</Link>.
        </>
    },
    2: {
        label: "Bloat",
        body: (cluster: ActiveCluster) => <OverviewBloat info={cluster}/>,
        info: <>
            You can reduce size of bloated tables and indexes without heavy locks here. It base on the
            tool <Link href={"https://github.com/dataegret/pgcompacttable"} target={"_blank"}>pgcompacttable</Link>.
            It is installed beside Ivory and will be used by Ivory. Ivory provides visualisation,
            helps to keep information about job and logs in one place by each cluster while you need it.
            Please be aware that you can run this tool only in master node and In target database contrib module
            pgstattuple should be installed via <i>create extension if not exists pgstattuple;</i>
        </>
    },
}

export type TabProps = {
    info: ActiveCluster
}

export function Overview() {
    const {activeCluster, activeClusterTab} = useStore()
    const {setClusterTab} = useStoreAction()
    const [infoOpen, setInfoOpen] = useState(false)
    const [settingsOpen, setSettingsOpen] = useState(false)
    const clusters = useQuery<ClusterMap>({queryKey: ["cluster/list"]})
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
        return tab.body(activeCluster)
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
