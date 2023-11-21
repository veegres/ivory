import {Box, IconButton, Tooltip} from "@mui/material";
import {Logout, Settings} from "@mui/icons-material";
import {useState} from "react";
import {randomUnicodeAnimal} from "../../app/utils";
import {useStoreAction} from "../../provider/StoreProvider";
import {AuthType, SxPropsMap} from "../../type/common";
import {useAuth} from "../../provider/AuthProvider";
import {useQueryClient} from "@tanstack/react-query";

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
    auth: AuthType,
    show: boolean,
}

export function Header(props: Props) {
    const {company, show, auth} = props
    const {toggleSettingsDialog} = useStoreAction()
    const {logout} = useAuth()
    const [animal, setAnimal] = useState("")

    const queryClient = useQueryClient();

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
            <>
                <Tooltip title={"Settings"}>
                    <IconButton onClick={toggleSettingsDialog} color={"inherit"}>
                        <Settings/>
                    </IconButton>
                </Tooltip>
                {auth === AuthType.BASIC && (
                    <Tooltip title={"Sign out"}>
                        <IconButton onClick={handleLogout} color={"inherit"}>
                            <Logout/>
                        </IconButton>
                    </Tooltip>
                )}
            </>
        )
    }

    function handleAnimal() {
        if (animal) setAnimal("")
        else setAnimal(randomUnicodeAnimal())
    }

    function handleLogout() {
        logout()
        queryClient.refetchQueries({queryKey: ["info"]}).then()
    }
}
