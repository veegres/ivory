import React from 'react';
import {NodeClusterInfoComponent} from "./features/node/NodeClusterInfoComponent";
import {Grid} from "@mui/material";
import {ClusterListComponent} from "./features/cluster/ClusterListComponent";
import {NodeConfigComponent} from "./features/node/NodeConfigComponent";
import {NodePatroniComponent} from "./features/node/NodePatroniComponent";
import {Item} from "./features/view/Item";
import {Header} from "./features/view/Header";

export function App() {
    return (
        <Grid container direction="column">
            <Grid item><Header /></Grid>
            <Grid item container><Item><ClusterListComponent /></Item></Grid>
            <Grid item container><Item><NodeClusterInfoComponent /></Item></Grid>
            <Grid item container><Item><NodePatroniComponent /></Item></Grid>
            <Grid item container><Item><NodeConfigComponent /></Item></Grid>
        </Grid>
    );
}
