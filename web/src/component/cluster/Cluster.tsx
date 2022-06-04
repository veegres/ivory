import {ClusterOverview} from "./ClusterOverview";
import {ClusterConfig} from "./ClusterConfig";
import {Alert, Box, Chip, Collapse, IconButton, Link, Tab, Tabs} from "@mui/material";
import React, {useState} from "react";
import {useStore} from "../../provider/StoreProvider";
import {ClusterBloat} from "./ClusterBloat";
import {Info} from "../view/Info";
import {Block} from "../view/Block";
import {useQuery} from "react-query";
import {InfoOutlined} from "@mui/icons-material";
import {ClusterTabs} from "../../app/types";

const SX = {
    headBox: {display: "flex", justifyContent: "space-between", alignItems: "center"},
    infoBox: {padding: '5px 0'},
    chip: {margin: "auto 0", minWidth: "150px"},
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
    const [isInfoOpen, setIsInfoOpen] = useState(false)
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
                <Box>
                    {!activeCluster.leader ? null : (
                        <Chip sx={SX.chip} color={"success"} label={activeCluster.leader} variant={"outlined"}/>
                    )}
                    {!tab?.info ? null : (
                        <IconButton color={"primary"} onClick={() => setIsInfoOpen(!isInfoOpen)}>
                            <InfoOutlined/>
                        </IconButton>
                    )}
                </Box>
            </Box>
            <Box sx={SX.infoBox}>{renderInfoBlock()}</Box>
            <Box>{renderMainBlock()}</Box>
        </Block>
    )

    function renderMainBlock() {
        if (disabled) return <Info text={"Please, select a cluster to see the information!"}/>
        if (!tab) return <Info text={"Coming soon — we're working on it!"}/>
        else return tab.body
    }

    function renderInfoBlock() {
        if (!tab?.info) return null
        return (
            <Collapse in={isInfoOpen}>
                <Alert severity="info" onClose={() => setIsInfoOpen(false)}>{tab.info}</Alert>
            </Collapse>
        )
    }

    function handleChange(value: number) {
        setStore({activeCluster: {...activeCluster, tab: value}})
    }
}
