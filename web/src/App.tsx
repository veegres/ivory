import React from 'react';
import {Cluster} from "./component/node/Cluster";
import {Grid} from "@mui/material";
import {Clusters} from "./component/clusters/Clusters";
import {Item} from "./component/view/Item";
import {Header} from "./component/view/Header";

export function App() {
    return (
        <Grid container direction="column">
            <Grid item><Header /></Grid>
            <Grid item container><Item><Clusters /></Item></Grid>
            <Grid item><Cluster /></Grid>
        </Grid>
    );
}
