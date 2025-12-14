import {Paper} from "@mui/material"

import {useRouterInfo} from "../../../../api/management/hook"
import {ErrorSmart} from "../../../view/box/ErrorSmart"
import {NoBox} from "../../../view/box/NoBox"
import {PermissionsList} from "./PermissionsList"

export function PermissionsSelfList() {
    const info = useRouterInfo()
    if (!info.data?.auth.user) return <ErrorSmart error={"Can't get user info"}/>
    const {username, permissions} = info.data.auth.user
    if (!permissions || Object.keys(permissions).length === 0) return <NoBox text={"You don't have any permissions"}/>
    return (
        <Paper variant={"outlined"}>
            <PermissionsList permissions={permissions} username={username}/>
        </Paper>
    )
}