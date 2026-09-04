import {Box} from "@mui/material"

import {ErrorSmart} from "../../../shared/component/box/ErrorSmart"
import {NoBox} from "../../../shared/component/box/NoBox"
import {SkeletonGroup} from "../../../shared/component/progress/SkeletonGroup"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {useRouterUserList} from "../api/UserHook"
import {UserListItem} from "./UserListItem"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
}

export function UserList() {
    const users = useRouterUserList()

    if (users.isLoading) return <SkeletonGroup count={3} width={"100%"} height={76} grow={false}/>
    if (users.isError) return <ErrorSmart error={users.error}/>
    if (!users.data || users.data.length === 0) return <NoBox text={"There are no users yet"}/>

    return (
        <Box sx={SX.box}>
            {users.data.map((user) => <UserListItem key={user.username} user={user}/>)}
        </Box>
    )
}
