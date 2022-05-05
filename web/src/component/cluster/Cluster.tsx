import {ClusterOverview} from "./ClusterOverview";
import {ClusterConfig} from "./ClusterConfig";
import {Box, FormControl, MenuItem, Select, Tab, Tabs} from "@mui/material";
import React, {useState} from "react";
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
    const {store: {activeCluster, activeNode}} = useStore()
    const disabled = !activeCluster.name

    const [tab, setTab] = useState(0)

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
            <Tabs value={tab} onChange={(_, value) => setTab(value)}>
                <Tab label={"Overview"} disabled={disabled} />
                <Tab label={"Config"} disabled={disabled}/>
                <Tab label={"Bloat"} disabled={disabled}/>
            </Tabs>
        )
    }

    function ActiveBlock() {
        switch (tab) {
            case 0:
                return <ClusterOverview cluster={activeCluster.name}/>
            case 1:
                return <ClusterConfig node={activeNode}/>
            case 2:
                return <ClusterBloat node={activeNode}/>
            default:
                return <Info text={"Coming soon — we're working on it!"}/>
        }
    }
}
