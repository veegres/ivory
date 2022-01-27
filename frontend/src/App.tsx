import React, {useState} from 'react';
import {NodeCluster} from "./features/node/NodeCluster";
import {Grid} from "@mui/material";
import {ClusterList} from "./features/cluster/ClusterList";
import {NodeConfig} from "./features/node/NodeConfig";
import {NodePatroni} from "./features/node/NodePatroni";
import {Item} from "./features/view/Item";
import {Header} from "./features/view/Header";

export function App() {
    const [node, setNode] = useState('')

    return (
        <Grid container direction="column">
            <Grid item><Header /></Grid>
            <Grid item container><Item><ClusterList setNode={setNode} /></Item></Grid>
            <Grid item container>
                <Item><NodePatroni node={node} /></Item>
                <Item><NodeCluster node={node} /></Item>
                <Item><NodeConfig node={node} /></Item>
            </Grid>
        </Grid>
    );
}
