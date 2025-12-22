import {Alert, Box, Collapse, Divider, Link, Tab, Tabs} from "@mui/material"
import {useMemo, useState} from "react"

import {useRouterClusterList, useRouterClusterOverview} from "../../../../api/cluster/hook"
import {ClusterTab} from "../../../../api/cluster/type"
import {useRouterInfo} from "../../../../api/management/hook"
import {Permission, PermissionStatus} from "../../../../api/permission/type"
import {SxPropsMap} from "../../../../app/type"
import {useStore, useStoreAction} from "../../../../provider/StoreProvider"
import {AlertCentered} from "../../../view/box/AlertCentered"
import {ErrorSmart} from "../../../view/box/ErrorSmart"
import {PageMainBox} from "../../../view/box/PageMainBox"
import {OverviewAction} from "./OverviewAction"
import {OverviewBloat} from "./OverviewBloat"
import {OverviewConfig} from "./OverviewConfig"
import {ClusterNoInstanceError} from "./OverviewError"
import {OverviewInstances} from "./OverviewInstances"
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

const TABS: ClusterTab[] = [
    {
        label: "Overview",
        permission: Permission.ViewInstanceOverview,
        body: (cluster, overview) => <OverviewInstances cluster={cluster} instances={overview?.instances}/>,
        info: <>
            The Overview tab offers visibility into the current status of your cluster. From here, you can
            utilize essential features to manage your cluster, including switchover, reinit, restart, reload,
            failover, and more. The leader node is automatically detected by sending requests to each instance
            until a successful connection is established. You have the flexibility to change the main instance
            to which Ivory sends requests by accessing the settings in the top right corner.
        </>
    },
    {
        label: "Config",
        permission: Permission.ViewInstanceConfig,
        body: (cluster, overview) => {
            if (!overview?.mainInstance) return <ClusterNoInstanceError/>
            return <OverviewConfig cluster={cluster} instance={overview.mainInstance}/>
        },
        info: <>
            You can adjust your PostgreSQL configurations here, and any changes made will be applied to
            all cluster instances. Instead of rewriting the entire configuration, it applies a patch
            update. If you wish to remove a specific setting, simply set it to <b>null</b>. Keep in mind that
            modifying certain parameters may necessitate restarting PostgreSQL. For further details on how
            this process functions, refer to
            the <Link href={"https://patroni.readthedocs.io/en/latest/rest_api.html#config-endpoint"} target={"_blank"}>documentation</Link>.
        </>
    },
    {
        label: "Bloat",
        permission: Permission.ViewBloatList,
        body: (cluster, overview) => {
            if (!overview?.mainInstance) return <ClusterNoInstanceError/>
            return <OverviewBloat cluster={cluster} instance={overview.mainInstance}/>
        },
        info: <>
            Here, you can efficiently decrease the size of bloated tables and indexes without imposing
            heavy locks. This functionality is powered by
            the <Link href={"https://github.com/dataegret/pgcompacttable"} target={"_blank"}>pgcompacttable</Link> tool,
            which is seamlessly integrated with Ivory for streamlined usage. Ivory simplifies visualization
            and centralizes information about jobs and logs within each cluster, ensuring convenient access when
            needed. It's important to note that this tool can only be executed on the master node, and
            in the target database, the contrib module pgstattuple must be installed using the command
            "<b>CREATE EXTENSION IF NOT EXISTS pgstattuple;</b>". Ivory supports such features as:
            <ul>
                <li><b>Delay ratio</b> - A dynamic part of the delay between rounds is calculated as previous-round-time * delay-ratio. By default 2.</li>
                <li><b>Min table size</b> - Tables smaller than the specified size (in megabytes) will be excluded from processing.</li>
                <li><b>Max table size</b> - Tables larger than the specified size (in megabytes) will be excluded from processing.</li>
                <li><b>Force</b> - Try to compact even those tables and indexes that do not meet minimal bloat requirements.</li>
                <li><b>Routing vacuum</b> - Turn on the routine vacuum. By default all the vacuums are off.</li>
                <li><b>No initial vacuum</b> - Turn off initial vacuum before table processing.</li>
                <li><b>Initial reindex</b> - Perform an initial reindex of tables before processing.</li>
                <li><b>No reindex</b> - Turn off reindexing of tables after processing.</li>
            </ul>
        </>
    },
]

export function Overview() {
    const activeClusterTab = useStore(s => s.activeClusterTab)
    const activeTags = useStore(s => s.activeTags)
    const activeCluster = useStore(s => s.activeCluster)
    const {setClusterTab} = useStoreAction
    const [infoOpen, setInfoOpen] = useState(false)
    const [settingsOpen, setSettingsOpen] = useState(false)
    const clusters = useRouterClusterList(activeTags)
    const overview = useRouterClusterOverview(activeCluster?.cluster.name, false)
    const info = useRouterInfo(false)
    const permissions = info.data?.auth.user?.permissions
    const tab = TABS[activeClusterTab]

    const visible = useMemo(handleMemoVisibility, [clusters.data, clusters.error])

    return (
        <PageMainBox withPadding visible={visible}>
            <Box sx={SX.headBox}>
                <Tabs value={activeClusterTab} onChange={(_, value) => setClusterTab(value)} role={"tab"}>
                    {TABS.map((value, i) => permissions && permissions[value.permission] === PermissionStatus.GRANTED && (
                        <Tab key={i} value={i} label={value.label}/>
                    ))}
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
        if (!activeCluster) return <AlertCentered text={"Please, select a cluster to see the overview! (click on the name)"}/>
        if (!tab) return <AlertCentered text={"Coming soon — we're working on it!"}/>
        return <>
            {overview.error && <ErrorSmart error={overview.error}/>}
            {tab.body(activeCluster.cluster, overview.data)}
        </>
    }

    function renderActions() {
        if (!activeCluster) return null
        return <OverviewAction
            cluster={activeCluster}
            mainInstance={overview.data?.mainInstance}
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
                    <OverviewOptions info={activeCluster} overview={overview.data}/>
                </Box>
            </Collapse>
        )
    }

    function handleMemoVisibility() {
        return Object.keys(clusters.data ?? {}).length !== 0 && clusters.error === null
    }
}
