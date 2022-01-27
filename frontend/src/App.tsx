import React, {useState} from 'react';
import {Node} from "./features/node/Node";
import {Grid} from "@mui/material";
import {ClusterList} from "./features/cluster/ClusterList";
import {Item} from "./features/view/Item";
import {Header} from "./features/view/Header";

export function App() {
    const [node, setNode] = useState('')

    return (
        <Grid container direction="column">
            <Grid item><Header /></Grid>
            <Grid item container><Item><ClusterList setNode={setNode} /></Item></Grid>
            <Grid item><Node name={node} /></Grid>
        </Grid>
    );
}
