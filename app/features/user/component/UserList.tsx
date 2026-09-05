import { Stars } from "@mui/icons-material"
import {Box, Tooltip} from "@mui/material"

import {ErrorSmart} from "../../../shared/component/box/ErrorSmart"
import {NoBox} from "../../../shared/component/box/NoBox"
import {TitleBox} from "../../../shared/component/box/TitleBox"
import {SkeletonGroup} from "../../../shared/component/progress/SkeletonGroup"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {useRouterPermissions} from "../../permission/api/PermissionHook"
import {PermissionList} from "../../permission/component/PermissionList"
import {useRouterUserList} from "../api/UserHook"
import {User} from "../api/UserType"
import {UserBasicStatus} from "./UserBasicStatus"
import {UserUpdate} from "./UserUpdate"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    body: {display: "flex", flexDirection: "column", gap: 1, marginTop: 0.5},
    actions: {display: "flex", alignItems: "center", gap: 0.5, justifyContent: "end", flexGrow: 1},
    perm: {maxHeight: "370px", overflow: "auto"},
    icon: {fontSize: "19px"},
}

export function UserList() {
    const permissions = useRouterPermissions()
    const users = useRouterUserList()

    if (permissions.isLoading || users.isLoading) return <SkeletonGroup count={3} width={"100%"} height={76} grow={false}/>
    if (permissions.isError) return <ErrorSmart error={permissions.error}/>
    if (users.isError) return <ErrorSmart error={users.error}/>
    if (!users.data || users.data.length === 0) return <NoBox text={"There are no users yet"}/>
    if (!permissions.data) return <NoBox text={"There is no permissions data"}/>

    return (
        <Box sx={SX.box}>
            {users.data.map((user) => (
                <TitleBox key={user.username}  label={user.username} renderActions={renderActions(user)} island={true}>
                    <Box sx={SX.body}>
                        <UserUpdate user={user}/>
                        {renderPermissions(user.username)}
                    </Box>
                </TitleBox>
            ))}
        </Box>
    )

    function renderPermissions(username: string) {
        const userPermissions = permissions.data?.find(e => e.username === username)?.permissions
        return (
            <Box sx={SX.perm}>
                <PermissionList permissions={userPermissions} username={username} view={"admin"}/>
            </Box>
        )
    }

    function renderActions(user: User) {
        return (
            <Box sx={SX.actions}>
                <UserBasicStatus status={user.registration?.status} expiresAt={user.registration?.expiresAt}/>
                <Tooltip title={user.superuser ? "Superuser" : "User"} placement={"top"} arrow disableInteractive>
                    <Stars sx={SX.icon} color={user.superuser ? "inherit" : "disabled"}/>
                </Tooltip>
            </Box>
        )
    }
}
