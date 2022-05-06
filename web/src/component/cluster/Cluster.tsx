import {ClusterOverview} from "./ClusterOverview";
import {ClusterConfig} from "./ClusterConfig";
import {Box, FormControl, MenuItem, Select, Tab, Tabs} from "@mui/material";
import React from "react";
import {useStore} from "../../provider/StoreProvider";
import {ClusterBloat} from "./ClusterBloat";
import {Info} from "../view/Info";
import {Block} from "../view/Block";


const SX = {
    headBox: {display: "flex", justifyContent: "space-between"},
    mainBox: {marginTop: '10px'},
    from: {minWidth: "170px", justifyContent: "center"},
}

export function Cluster() {
    const {store: {activeCluster}, setStore} = useStore()
    const disabled = !activeCluster.name

    // TODO add menu items, think how to show leader
    return (
        <Block withPadding>
            <Box sx={SX.headBox}>
                <Tabulation/>
                {!activeCluster.node ? null : (
                    <FormControl sx={SX.from} size={"small"} variant={"standard"}>
                        <Select value={activeCluster.node}>
                            <MenuItem value={activeCluster.node}>{activeCluster.node}</MenuItem>
                        </Select>
                    </FormControl>
                )}
            </Box>
            <Box sx={SX.mainBox}>
                {disabled ? (
                    <Info text={"Please, select a cluster to see the information!"} />
                ) : (
                    <ActiveBlock/>
                )}
            </Box>
        </Block>
    )

    function Tabulation() {
        return (
            <Tabs value={activeCluster.tab} onChange={(_, value) => setStore({activeCluster: {...activeCluster, tab: value}})}>
                <Tab label={"Overview"} disabled={disabled} />
                <Tab label={"Config"} disabled={disabled}/>
                <Tab label={"Bloat"} disabled={disabled}/>
            </Tabs>
        )
    }

    function ActiveBlock() {
        switch (activeCluster.tab) {
            case 0:
                return <ClusterOverview cluster={activeCluster.name}/>
            case 1:
                return <ClusterConfig node={activeCluster.node}/>
            case 2:
                return <ClusterBloat node={activeCluster.node}/>
            default:
                return <Info text={"Coming soon — we're working on it!"}/>
        }
    }
}
