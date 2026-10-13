import {Logout, PersonOutlined} from "@mui/icons-material"
import {Box, IconButton, Tooltip} from "@mui/material"
import {useState} from "react"

import {Settings} from "../../../core/widgets/settings/Settings"
import {useRouterLogout} from "../../../features/auth/api/AuthHook"
import {SxPropsMap} from "../../helper/HelperType"
import {randomUnicodeAnimal} from "../../helper/HelperUtils"

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "space-between", alignItems: "center", gap: {xs: 1, sm: 2, md: 8}, paddingX: {xs: 1, sm: "20px"}},
    main: {display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center", flexShrink: 0},
    title: {fontSize: {xs: "26px", sm: "35px"}, fontWeight: 900, fontFamily: "monospace", cursor: "pointer", color: "primary.main", userSelect: "none", whiteSpace: "nowrap"},
    caption: {display: {xs: "none", sm: "block"}, fontSize: "12px", fontWeight: 500, fontFamily: "monospace"},
    emblem: {padding: {xs: "5px 0", sm: "5px 8px"}, fontWeight: "bold", fontSize: {xs: "16px", sm: "20px"}, whiteSpace: "nowrap", textOverflow: "ellipsis", overflow: "hidden"},
    side: {display: "flex", flex: "1 1 0", alignItems: "center", gap: {xs: 0, sm: 1}, minWidth: 0},
    username: {display: {xs: "none", sm: "block"}, color: "text.secondary", padding: "5px", whiteSpace: "nowrap", textOverflow: "ellipsis", overflow: "hidden"},
    user: {display: {xs: "inline-flex", sm: "none"}},
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
    const [animal, setAnimal] = useState("")

    const logoutRouter = useRouterLogout()

    return (
        <Box sx={SX.box} >
            <Box sx={[SX.side, {justifyContent: "start"}]}>
               {show && (<Box sx={SX.emblem}>{company.toUpperCase()}</Box>)}
            </Box>
            <Box sx={SX.main}>
                <Box sx={SX.title} onClick={handleAnimal}>
                    {animal} Ivory {animal}
                </Box>
                <Box sx={SX.caption}>[database cluster management]</Box>
            </Box>
            <Box sx={[SX.side, {justifyContent: "end"}]}>
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
                {showLogout && (
                    <Tooltip title={username} enterTouchDelay={0} leaveTouchDelay={3000}>
                        <IconButton sx={SX.user} color={"inherit"}>
                            <PersonOutlined sx={SX.icon}/>
                        </IconButton>
                    </Tooltip>
                )}
                <Settings/>
                {showLogout && (
                    <Tooltip title={"Sign out"}>
                        <IconButton onClick={handleLogout} color={"inherit"}>
                            <Logout sx={SX.icon}/>
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
        logoutRouter.mutate()
    }
}
