import {Box} from "@mui/material"
import {useState} from "react"

import {PermissionAccountList} from "../../../features/permission/component/PermissionAccountList"
import {UserAccountPassword} from "../../../features/user/component/UserAccountPassword"
import {Tabs, TabsButton} from "../../../shared/component/button/TabsButton"
import {LastElementScrolling} from "../../../shared/component/scrolling/LastElementScrolling"

const TABS: Tabs = {
    0: {label: "Password"},
    1: {label: "Permissions"},
}

export function SettingsAccount() {
    const [tab, setTab] = useState(0)
    return (
        <LastElementScrolling>
            <TabsButton tabs={TABS} tab={tab} setTab={setTab}/>

            <Box>
                {tab === 0 && <UserAccountPassword/>}
                {tab === 1 && <PermissionAccountList/>}
            </Box>
        </LastElementScrolling>
    )
}