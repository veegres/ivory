import React from 'react';
import {Node} from "./component/node/Node";
import {Grid} from "@mui/material";
import {ClusterList} from "./component/cluster/ClusterList";
import {Item} from "./component/view/Item";
import {Header} from "./component/view/Header";

// TODO 1 Refactor all nested function components <> to render() or extract to separate file (diff behave could cause a problems)
export function App() {
    return (
        <Grid container direction="column">
            <Grid item><Header /></Grid>
            <Grid item container><Item><ClusterList /></Item></Grid>
            <Grid item><Node /></Grid>
        </Grid>
    );
}
