import {ErrorUserInfoMissing} from "../../../shared/component/box/ErrorManual"
import {NoBox} from "../../../shared/component/box/NoBox"
import {useRouterInfo} from "../../management/api/ManagementHook"
import {PermissionList} from "./PermissionList"

export function PermissionAccountList() {
    const info = useRouterInfo(false)
    if (!info.data?.auth.user) return <ErrorUserInfoMissing/>
    const {username, permissions} = info.data.auth.user
    if (!permissions || Object.keys(permissions).length === 0) return <NoBox text={"You don't have any permissions"}/>
    return (
        <PermissionList permissions={permissions} username={username}/>
    )
}