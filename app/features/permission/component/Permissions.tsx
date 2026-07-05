import {Box, Tab, Tabs} from "@mui/material"
import {useState} from "react"

import {LastElementScrolling} from "../../../shared/component/scrolling/LastElementScrolling"
import {Feature} from "../../Feature"
import {ManageAccess} from "../../management/component/ManageAccess"
import {PermissionsSelfList} from "./PermissionsSelfList"
import {PermissionsUserList} from "./PermissionsUserList"

export function Permissions() {
    const [tab, setTab] = useState<"self" | "users">("self")
    return (
        <LastElementScrolling>
            <ManageAccess feature={Feature.ViewPermissionList}>
                <Tabs variant={"fullWidth"} value={tab} onChange={(_, value) => setTab(value as "self" | "users")}>
                    <Tab label={"Self"} value={"self"}/>
                    <Tab label={"Users"} value={"users"}/>
                </Tabs>
            </ManageAccess>

            <Box>
                {tab === "self" && <PermissionsSelfList/>}
                {tab === "users" && <PermissionsUserList/>}
            </Box>
        </LastElementScrolling>
    )
}