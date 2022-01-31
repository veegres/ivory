import {Box, Grid, IconButton} from "@mui/material";
import {Brightness4, Brightness7} from "@mui/icons-material";
import React from "react";
import {useTheme} from "../../provider/ThemeProvider";

const SX = {
    title: {fontSize: '35px', fontWeight: 900, fontFamily: 'monospace'},
    sides: {fontSize: '20px', margin: '5px 20px', height: '100%', fontWeight: 900, borderRadius: '15px'},
    caption: {fontSize: '12px', fontWeight: 500, fontFamily: 'monospace'},
    container: {marginBottom: '10px'},
    emblem: {padding: '5px 15px'},
    buttons: {padding: '0px 5px'}
}

export function Header() {
    const theme = useTheme()
    const color = theme.info?.palette.primary.main

    return (
        <Grid container sx={SX.container} justifyContent="space-between" alignItems="center">
            <Grid item sx={{...SX.sides, borderLeft: `1px solid ${color}`}}>
                <Box sx={SX.emblem}>VEEGRES</Box>
            </Grid>
            <Grid item textAlign={"center"}>
                <Box sx={{...SX.title, color}}>Ivory</Box>
                <Box sx={SX.caption}>[postgres cluster visualization]</Box>
            </Grid>
            <Grid item sx={{...SX.sides, borderRight: `1px solid ${color}`}}>
                <Box sx={SX.buttons}>
                    <IconButton onClick={theme.toggle}>
                        {theme.mode === 'dark' ? <Brightness7/> : <Brightness4/>}
                    </IconButton>
                </Box>
            </Grid>
        </Grid>
    )
}
