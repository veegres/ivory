import {Box, Grid, IconButton} from "@mui/material";
import {Brightness4, Brightness7} from "@mui/icons-material";
import React from "react";
import {useTheme} from "../../provider/ThemeProvider";

const SX = {
    titleText: { fontSize: '35px', fontWeight: 900, fontFamily: 'monospace' },
    sidesText: { fontSize: '20px', padding: '0px 20px', fontWeight: 900 },
    captionText: { fontSize: '12px', fontWeight: 500, fontFamily: 'monospace' },
    container: { marginBottom: '10px' },
}

export function Header() {
    const theme = useTheme()
    const color = theme.info?.palette.primary.main

    return (
        <Grid container sx={SX.container} justifyContent="space-between" alignItems="center">
            <Grid item sx={SX.sidesText}>VEEGRES</Grid>
            <Grid item textAlign={"center"}>
                <Box sx={{ ...SX.titleText, color }}>Ivory</Box>
                <Box sx={SX.captionText}>[postgres cluster visualization]</Box>
            </Grid>
            <Grid item>
                <IconButton onClick={theme.toggle}>
                    {theme.mode === 'dark' ? <Brightness7 /> : <Brightness4 />}
                </IconButton>
            </Grid>
        </Grid>
    )
}
