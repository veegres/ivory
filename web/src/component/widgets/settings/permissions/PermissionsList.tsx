import {Box} from "@mui/material"

import {PermissionStatus} from "../../../../api/permission/type"
import {PermissionsListItem} from "./PermissionsListItem"

type Props = {
    username: string,
    permissions: { [permission: string]: PermissionStatus },
    view?: "admin" | "user",
}

export function PermissionsList(props: Props) {
    const {permissions, username, view = "user"} = props
    return (
        <Box width={"100%"}>
            {Object.entries(permissions).map(([name, status]) => (
                <PermissionsListItem key={name} username={username} name={name} status={status} view={view}/>
            ))}
        </Box>
    )
}