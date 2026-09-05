import {useState} from "react"

import {Feature} from "../../../features/Feature"
import {UserList} from "../../../features/user/component/UserList"
import {UserRegistration} from "../../../features/user/component/UserRegistration"
import {NoBox} from "../../../shared/component/box/NoBox"
import {Tabs, TabsButton} from "../../../shared/component/button/TabsButton"
import {LastElementScrolling} from "../../../shared/component/scrolling/LastElementScrolling"


const TABS: Tabs = {
    0: {label: "USERS", feature: Feature.ViewUserList},
    1: {label: "REGISTRATION", feature: Feature.ManageUserCreate},
}

export function SettingsUser() {
    const [tab, setTab] = useState(0)

    return (
        <LastElementScrolling>
            <TabsButton tabs={TABS} tab={tab} setTab={setTab}/>
            {renderTab()}
        </LastElementScrolling>
    )

    function renderTab() {
        switch (tab) {
            case 0: return <UserList/>
            case 1: return <UserRegistration/>
            default: return <NoBox text={"Tab is not supported"}/>
        }
    }
}
