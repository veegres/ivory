import React from 'react';
import {ClusterInfoComponent} from "./features/cluster/ClusterInfoComponent";
import {Box, Grid} from "@material-ui/core";
import {ClusterListComponent} from "./features/cluster/ClusterListComponent";
import {NodeConfigComponent} from "./features/node/NodeConfigComponent";
import {NodePatroniComponent} from "./features/node/NodePatroniComponent";

export default function App() {
    return (
        <Grid container spacing={3} direction="column">
            <Grid item><Box style={{textAlign: "center", fontSize: "30px"}}>Postgres Cluster GUI</Box></Grid>
            <Grid item><ClusterListComponent /></Grid>
            <Grid item><ClusterInfoComponent /></Grid>
            <Grid item><NodePatroniComponent /></Grid>
            <Grid item><NodeConfigComponent /></Grid>
        </Grid>
    );
}
