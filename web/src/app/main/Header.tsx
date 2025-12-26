import {Logout, Settings} from "@mui/icons-material"
import {Box, IconButton, Tooltip} from "@mui/material"
import {useState} from "react"

import {useRouterLogout} from "../../api/auth/hook"
import {ClusterApi} from "../../api/cluster/router"
import {Refresher} from "../../component/widgets/refresher/Refresher"
import {useStoreAction} from "../../provider/StoreProvider"
import {SxPropsMap} from "../type"
import {randomUnicodeAnimal} from "../utils"

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "space-between", alignItems: "center", gap: 8, padding: "0 20px"},
    main: {display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center"},
    title: {fontSize: "35px", fontWeight: 900, fontFamily: "monospace", cursor: "pointer", color: "primary.main", userSelect: "none"},
    caption: {fontSize: "12px", fontWeight: 500, fontFamily: "monospace"},
    emblem: {padding: "5px 8px", fontWeight: "bold", fontSize: "20px"},
    side: {display: "flex", flex: "1 1 0", alignItems: "center", gap: 1},
    username: {color: "text.secondary", padding: "5px"},
    icon: {fontSize: "22px"},
}

type Props = {
    company: string,
    username: string,
    show: boolean,
    showLogout: boolean,
}

export function Header(props: Props) {
    const {company, username, show, showLogout} = props
    const {toggleSettingsDialog} = useStoreAction
    const [animal, setAnimal] = useState("")

    const logoutRouter = useRouterLogout()

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
                {show && renderRightSide()}
            </Box>
        </Box>
    )

    function renderRightSide() {
        return (
            <>
                {showLogout && (
                    <Tooltip title={"Username"}>
                        <Box sx={SX.username}>{username}</Box>
                    </Tooltip>
                )}
                <Refresher queryKeys={[ClusterApi.list.key(), ClusterApi.overview.key()]}/>
                <Box>
                    <Tooltip title={"Settings"}>
                        <IconButton onClick={toggleSettingsDialog} color={"inherit"}>
                            <Settings sx={SX.icon}/>
                        </IconButton>
                    </Tooltip>
                    {showLogout && (
                        <Tooltip title={"Sign out"}>
                            <IconButton onClick={handleLogout} color={"inherit"}>
                                <Logout sx={SX.icon}/>
                            </IconButton>
                        </Tooltip>
                    )}
                </Box>
            </>
        )
    }

    function handleAnimal() {
        if (animal) setAnimal("")
        else setAnimal(randomUnicodeAnimal())
    }

    function handleLogout() {
        logoutRouter.mutate()
    }
}
