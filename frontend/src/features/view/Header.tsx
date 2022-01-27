import {Grid, IconButton} from "@mui/material";
import {Brightness4, Brightness7} from "@mui/icons-material";
import React from "react";
import {useTheme} from "../../provider/ThemeProvider";

export function Header() {
    const theme = useTheme()

    return (
        <Grid container justifyContent="space-between" alignItems="center">
            <Grid item sx={{ fontSize: '20px', padding: '0px 20px', fontWeight: 900 }}>VEEGRES</Grid>
            <Grid item sx={{ fontSize: '30px' }}>Postgres Cluster GUI</Grid>
            <Grid item>
                <IconButton onClick={theme.toggle}>
                    {theme.mode === 'dark' ? <Brightness7 /> : <Brightness4 />}
                </IconButton>
            </Grid>
        </Grid>
    )
}
