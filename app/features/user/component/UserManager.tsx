import {Box} from "@mui/material"
import {useState} from "react"

import {AlertCentered} from "../../../shared/component/box/AlertCentered"
import {Tabs, TabsButton} from "../../../shared/component/button/TabsButton"
import {LastElementScrolling} from "../../../shared/component/scrolling/LastElementScrolling"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {Feature} from "../../Feature"
import {useHasAccess} from "../../management/component/ManageAccess"
import {UserCreation} from "./UserCreation"
import {UserCreationForm} from "./UserCreationForm"
import {UserList} from "./UserList"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 2, padding: "5px 10px 0px 0px"},
}

const ACCOUNT_TAB = 0
const USERS_TAB = 1
const REGISTER_TAB = 2

const ACCOUNT_TABS: Tabs = {
    [ACCOUNT_TAB]: {label: "ACCOUNT"},
}

export function UserManager() {
    const [tab, setTab] = useState(ACCOUNT_TAB)
    const listAccess = useHasAccess(Feature.ViewUserList)
    const createAccess = useHasAccess(Feature.ManageUserCreate)

    const tabs: Tabs = {...ACCOUNT_TABS}
    if (listAccess === "allowed") tabs[USERS_TAB] = {label: "USERS"}
    if (createAccess === "allowed") tabs[REGISTER_TAB] = {label: "REGISTER"}

    return (
        <LastElementScrolling>
            {Object.keys(tabs).length > 1 && <TabsButton tabs={tabs} tab={tab} setTab={setTab}/>}
            <Box sx={SX.box}>{renderTab()}</Box>
        </LastElementScrolling>
    )

    function renderTab() {
        switch (tab) {
            case USERS_TAB: return <UserList/>
            case REGISTER_TAB: return renderRegister()
            default: return <UserCreation/>
        }
    }

    function renderRegister() {
        return (<>
            <AlertCentered text={renderDescription()}/>
            <UserCreationForm/>
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
