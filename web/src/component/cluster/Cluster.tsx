import {ClusterOverview} from "./ClusterOverview";
import {ClusterConfig} from "./ClusterConfig";
import {Box, Chip, Tab, Tabs} from "@mui/material";
import React from "react";
import {useStore} from "../../provider/StoreProvider";
import {ClusterBloat} from "./ClusterBloat";
import {Info} from "../view/Info";
import {Block} from "../view/Block";
import {useQuery} from "react-query";


const SX = {
    headBox: {display: "flex", justifyContent: "space-between"},
    mainBox: {marginTop: '10px'},
    chip: {margin: "auto 0", minWidth: "150px"},
}

export function Cluster() {
    const {store: {activeCluster}, setStore} = useStore()
    const disabled = !activeCluster.name

    const clusters = useQuery('cluster/list')

    return (
        <Block withPadding visible={clusters.isSuccess}>
            <Box sx={SX.headBox}>
                <Tabs value={activeCluster.tab} onChange={(_, value) => handleChange(value)}>
                    <Tab label={"Overview"} disabled={disabled} />
                    <Tab label={"Config"} disabled={disabled}/>
                    <Tab label={"Bloat"} disabled={disabled}/>
                </Tabs>
                {!activeCluster.leader ? null : (
                    <Chip sx={SX.chip} color={"success"} label={activeCluster.leader} variant={"outlined"} />
                )}
            </Box>
            <Box sx={SX.mainBox}>
                {!disabled ? renderActiveBlock() : (
                    <Info text={"Please, select a cluster to see the information!"} />
                )}
            </Box>
        </Block>
    )

    function renderActiveBlock() {
        switch (activeCluster.tab) {
            case 0:
                return <ClusterOverview cluster={activeCluster.name}/>
            case 1:
                return <ClusterConfig leader={activeCluster.leader}/>
            case 2:
                return <ClusterBloat leader={activeCluster.leader} cluster={activeCluster.name}/>
            default:
                return <Info text={"Coming soon — we're working on it!"}/>
        }
    }

    function handleChange(value: number) {
        setStore({activeCluster: {...activeCluster, tab: value}})
    }
}
