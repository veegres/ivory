import React from 'react';
import {NodeClusterInfoComponent} from "./features/node/NodeClusterInfoComponent";
import {Grid, IconButton} from "@mui/material";
import {ClusterListComponent} from "./features/cluster/ClusterListComponent";
import {NodeConfigComponent} from "./features/node/NodeConfigComponent";
import {NodePatroniComponent} from "./features/node/NodePatroniComponent";
import {useTheme} from "./provider/ThemeProvider";
import {Brightness4, Brightness7} from "@mui/icons-material";
import {Item} from "./features/view/Item";

export function App() {
    const theme = useTheme()
    return (
        <Grid container direction="column">
            <Grid item container justifyContent="space-between" alignItems="center">
                <Grid item sx={{ fontSize: '20px', padding: '0px 20px', fontWeight: 900 }}>VEEGRES</Grid>
                <Grid item sx={{ fontSize: '30px' }}>Postgres Cluster GUI</Grid>
                <Grid item>
                    <IconButton onClick={theme.toggle}>
                        {theme.mode === 'dark' ? <Brightness7 /> : <Brightness4 />}
                    </IconButton>
                </Grid>
            </Grid>
            <Grid item container><Item><ClusterListComponent /></Item></Grid>
            <Grid item container><Item><NodeClusterInfoComponent /></Item></Grid>
            <Grid item container><Item><NodePatroniComponent /></Item></Grid>
            <Grid item container><Item><NodeConfigComponent /></Item></Grid>
        </Grid>
    );
}
