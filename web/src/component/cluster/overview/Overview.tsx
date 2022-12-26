import {OverviewInstances} from "./OverviewInstances";
import {OverviewConfig} from "./OverviewConfig";
import {
    Alert, Box, Collapse, Divider, Link,
    Tab, Tabs, ToggleButton, ToggleButtonGroup, Tooltip
} from "@mui/material";
import React, {useState} from "react";
import {useStore} from "../../../provider/StoreProvider";
import {OverviewBloat} from "./OverviewBloat";
import {InfoAlert} from "../../view/InfoAlert";
import {Block} from "../../view/Block";
import {useQuery} from "@tanstack/react-query";
import {Article, InfoOutlined, Settings, Warning} from "@mui/icons-material";
import {ClusterTabs, CredentialType, ActiveCluster} from "../../../app/types";
import {OverviewSettings} from "./OverviewSettings";
import {InfoIcons} from "../../view/InfoIcons";
import {CredentialOptions, InstanceColor} from "../../../app/utils";
import {orange, purple} from "@mui/material/colors";
import {InfoBox} from "../../view/InfoBox";
import {InfoTitle} from "../../view/InfoTitle";

const SX = {
    headBox: {display: "flex", justifyContent: "space-between", alignItems: "center"},
    infoBox: {padding: '5px 0'},
    rightBox: {display: "flex", alignItems: "center", gap: 1},
    chip: {margin: "auto 0", borderRadius: "4px"},
    settingsBox: {height: "100%", display: "flex", flexDirection: "row"},
    mainBox: {display: "flex", flexDirection: "row"},
    leftMainBlock: {flexGrow: 1, overflowX: "auto"},
    rightMainBlock: {},
    dividerVertical: {margin: "0 10px"},
    dividerHorizontal: {margin: "10px 0"},
    collapse: {height: "100%"},
    toggleButton: {padding: "3px"}
}

const TABS: ClusterTabs = {
    0: {
        body: (cluster: ActiveCluster) => <OverviewInstances info={cluster}/>,
    },
    1: {
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
    const {store: {activeCluster, activeClusterTab}, setClusterTab} = useStore()
    const [infoOpen, setInfoOpen] = useState(false)
    const [settingsOpen, setSettingsOpen] = useState(false)
    const clusters = useQuery(["cluster/list"])
    const tab = TABS[activeClusterTab]

    return (
        <Block withPadding visible={clusters.isSuccess}>
            <Box sx={SX.headBox}>
                <Tabs value={activeClusterTab} onChange={(_, value) => handleChange(value)}>
                    <Tab label={"Instances"}/>
                    <Tab label={"Config"}/>
                    <Tab label={"Bloat"}/>
                </Tabs>
                {renderActionBlock()}
            </Box>
            <Box sx={SX.infoBox}>{renderInfoBlock()}</Box>
            <Box sx={SX.mainBox}>
                <Box sx={SX.leftMainBlock}>{renderMainBlock()}</Box>
                <Box sx={SX.rightMainBlock}>{renderSettingsBlock()}</Box>
            </Box>
        </Block>
    )

    function renderMainBlock() {
        if (!activeCluster) return <InfoAlert text={"Please, select a cluster to see the information!"}/>
        if (!tab) return <InfoAlert text={"Coming soon — we're working on it!"}/>
        return tab.body(activeCluster)
    }

    function renderActionBlock() {
        const cluster = activeCluster?.cluster

        return (
            <Box sx={SX.rightBox}>
                {renderShortClusterInfo()}
                <ToggleButtonGroup size={"small"}>
                    <ToggleButton
                        sx={SX.toggleButton}
                        value={"settings"}
                        disabled={!cluster}
                        selected={cluster && settingsOpen}
                        onClick={() => setSettingsOpen(!settingsOpen)}
                    >
                        <Tooltip title={"luster Settings"} placement={"top"}><Settings/></Tooltip>
                    </ToggleButton>
                    <ToggleButton
                        sx={SX.toggleButton}
                        value={"info"}
                        selected={!!tab?.info && infoOpen}
                        disabled={!tab?.info}
                        onClick={() => setInfoOpen(!infoOpen)}
                    >
                        <Tooltip title={"Tab Information"} placement={"top"}><InfoOutlined/></Tooltip>
                    </ToggleButton>
                </ToggleButtonGroup>
            </Box>
        )
    }

    function renderShortClusterInfo() {
        if (!activeCluster) return null
        const {cluster, instance, warning, detection} = activeCluster
        const postgres = CredentialOptions[CredentialType.POSTGRES]
        const patroni = CredentialOptions[CredentialType.PATRONI]

        const infoItems = [
            { icon: <Article />, name: "Patroni Certs", active: !!cluster.certId },
            { icon: patroni.icon, name: "Patroni Password", active: !!cluster.patroniCredId },
            { icon: postgres.icon, name: "Postgres Password", active: !!cluster.postgresCredId }
        ]
        const warningItems = [
            { icon: <Warning />,  name: "Warning", active: warning, iconColor: orange[500] }
        ]
        const roleTooltip = [
            { name: "Detection", value: detection, bgColor: purple[400] },
            { name: "Instance", value: instance.api_domain, bgColor: InstanceColor[instance.role] }
        ]

        return (
            <Box sx={SX.rightBox}>
                <InfoIcons items={warningItems} />
                <InfoIcons items={infoItems} />
                <InfoBox tooltip={<InfoTitle items={roleTooltip} />} withPadding>
                    <Box sx={{color: InstanceColor[instance.role]}}>
                        {instance.role.toUpperCase()}
                    </Box>
                </InfoBox>
            </Box>
        )
    }

    function renderInfoBlock() {
        if (!tab?.info) return null
        return (
            <Collapse in={infoOpen}>
                <Alert severity="info" onClose={() => setInfoOpen(false)}>{tab.info}</Alert>
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
                    <OverviewSettings info={activeCluster} />
                </Box>
            </Collapse>
        )
    }

    function handleChange(value: number) {
        setClusterTab(value)
    }
}
