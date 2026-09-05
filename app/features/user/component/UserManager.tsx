import {Box} from "@mui/material"
import {useState} from "react"

import {AlertCentered} from "../../../shared/component/box/AlertCentered"
import {Tabs, TabsButton} from "../../../shared/component/button/TabsButton"
import {LastElementScrolling} from "../../../shared/component/scrolling/LastElementScrolling"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {Feature} from "../../Feature"
import {UserAccount} from "./UserAccount"
import {UserRegistrationForm} from "./UserRegistrationForm"
import {UserList} from "./UserList"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 2, padding: "5px 10px 0px 0px"},
}

const TABS: Tabs = {
    0: {label: "ACCOUNT"},
    1: {label: "USERS", feature: Feature.ViewUserList},
    2: {label: "REGISTER", feature: Feature.ManageUserCreate},
}

export function UserManager() {
    const [tab, setTab] = useState(0)

    return (
        <LastElementScrolling>
            <TabsButton tabs={TABS} tab={tab} setTab={setTab}/>
            <Box sx={SX.box}>{renderTab()}</Box>
        </LastElementScrolling>
    )

    function renderTab() {
        switch (tab) {
            case 1: return <UserList/>
            case 2: return renderRegister()
            default: return <UserAccount/>
        }
    }

    function renderRegister() {
        return (<>
            <AlertCentered text={renderDescription()}/>
            <UserRegistrationForm/>
        </>)
    }

    function renderDescription() {
        return (
            "A user is registered by name and by the ways they may sign in - nobody who was not " +
            "registered gets in, whichever directory vouches for them. A user signing in with an " +
            "Ivory password sets it themselves on the page their one-time link opens."
        )
    }
}
