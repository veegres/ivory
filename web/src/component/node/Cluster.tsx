import {Item} from "../view/Item";
import {ClusterOverview} from "./ClusterOverview";
import {ClusterConfig} from "./ClusterConfig";
import {Alert, Box, FormControl, Grid, MenuItem, Select, Tab, Tabs} from "@mui/material";
import React, {useState} from "react";
import {useStore} from "../../provider/StoreProvider";
import {ClusterBloat} from "./ClusterBloat";


const SX = {
    headBox: {padding: "0 20px", display: "flex", justifyContent: "space-between"},
    mainBox: {padding: '10px 20px 15px'},
    from: {minWidth: "170px", justifyContent: "center"},
    infoAlert: {justifyContent: 'center'}
}

type InfoProps = { text: string }

export function Cluster() {
    const {store: {activeCluster, activeNode}} = useStore()
    const [tab, setTab] = useState(0)

    return (
        <Grid container>
            <Item>
                <Box sx={SX.headBox}>
                    <Tabulation/>
                    {!activeNode ? null : (
                        <FormControl sx={SX.from} size={"small"} variant={"standard"}>
                            <Select value={activeNode}>
                                <MenuItem value={activeNode}>{activeNode}</MenuItem>
                            </Select>
                        </FormControl>
                    )}
                </Box>
                <Box sx={SX.mainBox}>
                    {!activeCluster ? <NonSelectedBlock/> : <ActiveBlock/>}
                </Box>
            </Item>
        </Grid>
    )

    function Tabulation() {
        return (
            <Tabs value={tab} onChange={(_, value) => setTab(value)}>
                <Tab label={"Overview"} disabled={!activeCluster} />
                <Tab label={"Config"} disabled={!activeCluster}/>
                <Tab label={"Bloat"} disabled={!activeCluster}/>
            </Tabs>
        )
    }

    function NonSelectedBlock() {
        return <Info text={"Please, select a cluster to see the information!"} />
    }

    function ActiveBlock() {
        switch (tab) {
            case 0:
                return <ClusterOverview cluster={activeCluster}/>
            case 1:
                return <ClusterConfig node={activeNode}/>
            case 2:
                return <ClusterBloat node={activeNode}/>
            default:
                return <Info text={"Coming soon — we're working on it!"}/>
        }
    }

    function Info(props: InfoProps) {
        return (
            <Alert sx={SX.infoAlert} severity={"info"} variant={"outlined"} icon={false}>
                {props.text}
            </Alert>
        )
    }
}
