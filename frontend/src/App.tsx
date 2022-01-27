import React from 'react';
import {Node} from "./features/node/Node";
import {Grid} from "@mui/material";
import {ClusterList} from "./features/cluster/ClusterList";
import {Item} from "./features/view/Item";
import {Header} from "./features/view/Header";

export function App() {
    return (
        <Grid container direction="column">
            <Grid item><Header /></Grid>
            <Grid item container><Item><ClusterList /></Item></Grid>
            <Grid item><Node /></Grid>
        </Grid>
    );
}
