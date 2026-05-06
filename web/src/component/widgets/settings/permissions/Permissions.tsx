import {Box, Tab, Tabs} from "@mui/material"
import {useState} from "react"

import {Feature} from "../../../../api/feature"
import {Access} from "../../access/Access"
import {MenuWrapper} from "../menu/MenuWrapper"
import {PermissionsSelfList} from "./PermissionsSelfList"
import {PermissionsUserList} from "./PermissionsUserList"

export function Permissions() {
    const [tab, setTab] = useState<"self" | "users">("self")
    return (
        <MenuWrapper>
            <Access feature={Feature.ViewPermissionList}>
                <Tabs variant={"fullWidth"} value={tab} onChange={(_, value) => setTab(value as "self" | "users")}>
                    <Tab label={"Self"} value={"self"}/>
                    <Tab label={"Users"} value={"users"}/>
                </Tabs>
            </Access>

            <Box>
                {tab === "self" && <PermissionsSelfList/>}
                {tab === "users" && <PermissionsUserList/>}
            </Box>
        </MenuWrapper>
    )
}