import React from 'react';
import {Cluster} from "./component/cluster/Cluster";
import {Grid} from "@mui/material";
import {Clusters} from "./component/clusters/Clusters";
import {Header} from "./component/view/Header";
import {Node} from "./component/node/Node"

export function App() {
    return (
        <Grid container direction="column">
            <Grid><Header /></Grid>
            <Clusters />
            <Cluster />
            <Node />
        </Grid>
    );
}
