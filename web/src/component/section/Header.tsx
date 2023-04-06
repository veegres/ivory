import {Box, IconButton, Tooltip} from "@mui/material";
import {Settings} from "@mui/icons-material";
import {useState} from "react";
import {randomUnicodeAnimal} from "../../app/utils";
import {useStore} from "../../provider/StoreProvider";
import {SxPropsMap} from "../../type/common";

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "space-between", alignItems: "center", gap: 8, padding: "0 20px"},
    main: {display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center"},
    title: {fontSize: "35px", fontWeight: 900, fontFamily: "monospace", cursor: "pointer", color: "primary.main", userSelect: "none"},
    caption: {fontSize: "12px", fontWeight: 500, fontFamily: "monospace"},
    emblem: {padding: "5px 8px", fontWeight: "bold", fontSize: "20px"},
    side: {display: "flex", flex: "1 1 0"},
}

type Props = {
    company: string,
    show: boolean,
}

export function Header(props: Props) {
    const {company, show} = props
    const {toggleSettingsDialog} = useStore()
    const [animal, setAnimal] = useState("")

    return (
        <Box sx={SX.box} >
            <Box sx={SX.side} justifyContent={"start"}>
               {show && (<Box sx={SX.emblem}>{company.toUpperCase()}</Box>)}
            </Box>
            <Box sx={SX.main}>
                <Box sx={SX.title} onClick={handleAnimal}>
                    {animal} Ivory {animal}
                </Box>
                <Box sx={SX.caption}>[postgres cluster management]</Box>
            </Box>
            <Box sx={SX.side} justifyContent={"end"}>
                {show && (<Box>{renderRightSide()}</Box>)}
            </Box>
        </Box>
    )

    function renderRightSide() {
        return (
            <Tooltip title={"Settings"}>
                <IconButton onClick={toggleSettingsDialog} color={"inherit"}>
                    <Settings/>
                </IconButton>
            </Tooltip>
        )
    }

    function handleAnimal() {
        if (animal) setAnimal("")
        else setAnimal(randomUnicodeAnimal())
    }
}
