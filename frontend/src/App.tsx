import React from 'react';
import {NodeCluster} from "./features/node/NodeCluster";
import {Grid} from "@mui/material";
import {ClusterList} from "./features/cluster/ClusterList";
import {NodeConfig} from "./features/node/NodeConfig";
import {NodePatroni} from "./features/node/NodePatroni";
import {Item} from "./features/view/Item";
import {Header} from "./features/view/Header";

export function App() {
    return (
        <Grid container direction="column">
            <Grid item><Header /></Grid>
            <Grid item container><Item><ClusterList /></Item></Grid>
            <Grid item container><Item><NodeCluster /></Item></Grid>
            <Grid item container><Item><NodePatroni /></Item></Grid>
            <Grid item container><Item><NodeConfig /></Item></Grid>
        </Grid>
    );
}
