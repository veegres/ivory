import {Box, Grid, IconButton} from "@mui/material";
import {BadgeOutlined, LightMode, Nightlight, Security} from "@mui/icons-material";
import React, {useState} from "react";
import {useTheme} from "../../provider/ThemeProvider";
import {randomUnicodeAnimal} from "../../app/utils";
import {useStore} from "../../provider/StoreProvider";
import select from "../../style/select.module.css";

const SX = {
    title: {fontSize: '35px', fontWeight: 900, fontFamily: 'monospace', cursor: "pointer"},
    sides: {fontSize: '20px', margin: '5px 20px', height: '100%', fontWeight: 900, borderRadius: '15px', flex: "1 1 0"},
    caption: {fontSize: '12px', fontWeight: 500, fontFamily: 'monospace'},
    emblem: {padding: '5px 15px'},
    buttons: {padding: '0px 5px', float: "right"}
}

type Props = {
    company?: string,
    show: boolean,
}

export function Header(props: Props) {
    const { company, show } = props
    const theme = useTheme()
    const { toggleCredentialsWindow, toggleCertsWindow } = useStore()
    const [animal, setAnimal] = useState("")
    const color = theme.info?.palette.primary.main

    return (
        <Grid container justifyContent={"space-between"} alignItems={"center"} flexWrap={"nowrap"}>
            <Grid item sx={{...SX.sides, borderLeft: `1px solid ${color}`}}>
                {show && <Box sx={SX.emblem}>{company ? company.toUpperCase() : "VEEGRES"}</Box>}
            </Grid>
            <Grid item textAlign={"center"}>
                <Box sx={{...SX.title, color}} className={select.none} onClick={handleAnimal}>{animal} Ivory {animal}</Box>
                <Box sx={SX.caption}>[postgres cluster management]</Box>
            </Grid>
            <Grid item sx={{...SX.sides, borderRight: `1px solid ${color}`}}>
                {show && (
                    <Box sx={SX.buttons}>
                        <IconButton onClick={toggleCertsWindow}>
                            <BadgeOutlined/>
                        </IconButton>
                        <IconButton onClick={toggleCredentialsWindow}>
                            <Security/>
                        </IconButton>
                        <IconButton onClick={theme.toggle}>
                            {theme.mode === 'dark' ? <LightMode/> : <Nightlight/>}
                        </IconButton>
                    </Box>
                )}
            </Grid>
        </Grid>
    )

    function handleAnimal() {
        if (animal) setAnimal("")
        else setAnimal(randomUnicodeAnimal())
    }
}
