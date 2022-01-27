import {Item} from "../view/Item";
import {NodePatroni} from "./NodePatroni";
import {NodeCluster} from "./NodeCluster";
import {NodeConfig} from "./NodeConfig";
import {Alert, Box, Grid, Tab, Tabs} from "@mui/material";
import React, {useState} from "react";
import {useStore} from "../../provider/StoreProvider";

export function Node() {
    const { store } = useStore()
    const [tab, setTab] = useState(0)

    return (
        <Grid container>
            <Item>
                <Tabulation />
                <Box sx={{ padding: '10px 20px 15px' }}>
                    <CurrentBlock />
                </Box>
            </Item>
        </Grid>
    )

    function Tabulation() {
        return (
            <Tabs value={tab} onChange={(_, value) => setTab(value)} centered disabled={!store.activeNode}>
                <Tab label={"Overview"} />
                <Tab label={"Cluster"} />
                <Tab label={"Config"} />
                <Tab label={"Cleaning"} />
            </Tabs>
        )
    }

    function CurrentBlock() {
        switch (tab) {
            case 0: return <NodePatroni node={store.activeNode} />
            case 1: return <NodeCluster node={store.activeNode} />
            case 2: return <NodeConfig node={store.activeNode} />
            default: return <Alert severity={"info"}>Coming soon — we're working on it!</Alert>
        }
    }
}
