import {Paper} from "@mui/material"

import {ErrorUserInfoMissing} from "../../../shared/component/box/ErrorManual"
import {NoBox} from "../../../shared/component/box/NoBox"
import {useRouterInfo} from "../../management/api/ManagementHook"
import {PermissionsList} from "./PermissionsList"

export function PermissionsSelfList() {
    const info = useRouterInfo(false)
    if (!info.data?.auth.user) return <ErrorUserInfoMissing/>
    const {username, permissions} = info.data.auth.user
    if (!permissions || Object.keys(permissions).length === 0) return <NoBox text={"You don't have any permissions"}/>
    return (
        <Paper variant={"outlined"}>
            <PermissionsList permissions={permissions} username={username}/>
        </Paper>
    )
}