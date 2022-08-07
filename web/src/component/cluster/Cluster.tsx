import {ClusterOverview} from "./ClusterOverview";
import {ClusterConfig} from "./ClusterConfig";
import {
    Alert, Box, Chip, Collapse, Divider, Fade, Link,
    Tab, Tabs, ToggleButton, ToggleButtonGroup, Tooltip
} from "@mui/material";
import React, {useState} from "react";
import {useStore} from "../../provider/StoreProvider";
import {ClusterBloat} from "./ClusterBloat";
import {Info} from "../view/Info";
import {Block} from "../view/Block";
import {useQuery} from "react-query";
import {InfoOutlined, Settings} from "@mui/icons-material";
import {ClusterTabs} from "../../app/types";
import {ClusterSettings} from "./ClusterSettings";

const SX = {
    headBox: {display: "flex", justifyContent: "space-between", alignItems: "center"},
    infoBox: {padding: '5px 0'},
    actionBox: {display: "flex", alignItems: "center", gap: 1},
    chip: {margin: "auto 0", minWidth: "150px"},
    settingsBox: {height: "100%", display: "flex", flexDirection: "row"},
    mainBox: {display: "flex", flexDirection: "row"},
    leftMainBlock: {flexGrow: 1, overflowX: "auto"},
    rightMainBlock: {},
    divider: {margin: "0 10px"},
    collapse: {height: "100%"}
}

const TABS: ClusterTabs = {
    0: {
        body: <ClusterOverview/>,
    },
    1: {
        body: <ClusterConfig/>,
        info: <>
            You can change your postgres configurations here (it will be applied on all cluster nodes).
            It doesn't rewrite all your config, it call patches the existing configuration. If you want to
            remove (reset) some setting just patch it with <i>null</i>. Be aware that some of the parameters
            requires a restart of postgres. More information how it works you can find in
            patroni <Link href={"https://patroni.readthedocs.io/en/latest/SETTINGS.html"} target={"_blank"}>site</Link>.
        </>
    },
    2: {
        body: <ClusterBloat/>,
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

export function Cluster() {
    const {store: {activeCluster}, setStore} = useStore()
    const [infoOpen, setInfoOpen] = useState(false)
    const [settingsOpen, setSettingsOpen] = useState(false)
    const disabled = !activeCluster.name
    const clusters = useQuery('cluster/list')
    const tab = TABS[activeCluster.tab]

    return (
        <Block withPadding visible={clusters.isSuccess}>
            <Box sx={SX.headBox}>
                <Tabs value={activeCluster.tab} onChange={(_, value) => handleChange(value)}>
                    <Tab label={"Overview"}/>
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
        if (disabled) return <Info text={"Please, select a cluster to see the information!"}/>
        if (!tab) return <Info text={"Coming soon — we're working on it!"}/>
        return tab.body
    }

    function renderActionBlock() {
        return (
            <Box sx={SX.actionBox}>
                <Box>
                    <Fade in={!!activeCluster.leader}>
                        <Tooltip title={"All requests go to this node!"} placement={"top"}>
                            <Chip sx={SX.chip} label={activeCluster.leader?.api_domain} variant={"outlined"}/>
                        </Tooltip>
                    </Fade>
                </Box>
                <ToggleButtonGroup size={"small"}>
                    <ToggleButton
                        value={"settings"}
                        disabled={disabled}
                        selected={!disabled && settingsOpen}
                        onClick={() => setSettingsOpen(!settingsOpen)}
                    >
                        <Tooltip title={"Cluster Settings"} placement={"top"}><Settings/></Tooltip>
                    </ToggleButton>
                    <ToggleButton
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

    function renderInfoBlock() {
        if (!tab?.info) return null
        return (
            <Collapse in={infoOpen}>
                <Alert severity="info" onClose={() => setInfoOpen(false)}>{tab.info}</Alert>
            </Collapse>
        )
    }

    function renderSettingsBlock() {
        return (
            <Collapse sx={SX.collapse} in={!disabled && settingsOpen} orientation={"horizontal"}>
                <Box sx={SX.settingsBox}>
                    <Divider sx={SX.divider} orientation={"vertical"} flexItem/>
                    {!disabled && settingsOpen && <ClusterSettings />}
                </Box>
            </Collapse>
        )
    }

    function handleChange(value: number) {
        setStore({activeCluster: {...activeCluster, tab: value}})
    }
}
