import React from 'react';
import {Cluster} from "./component/cluster/Cluster";
import {Grid} from "@mui/material";
import {Clusters} from "./component/clusters/Clusters";
import {Block} from "./component/view/Block";
import {Header} from "./component/view/Header";
import {Node} from "./component/node/Node"

export function App() {
    return (
        <Grid container direction="column">
            <Grid item><Header /></Grid>
            <Grid item container><Block><Clusters /></Block></Grid>
            <Grid item container><Block withPadding><Cluster /></Block></Grid>
            <Grid item container><Block withPadding><Node /></Block></Grid>
        </Grid>
    );
}
