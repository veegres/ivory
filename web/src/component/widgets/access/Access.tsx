import {Box, BoxProps} from "@mui/material"
import {cloneElement} from "react"

import {useRouterInfo} from "../../../api/management/hook"
import {Permission, PermissionStatus} from "../../../api/permission/type"

type Props = BoxProps & {
    permission: Permission,
}

export function Access(props: Props) {
    const info = useRouterInfo()
    const permissions = info.data?.auth.user?.permissions
    if (permissions && permissions[props.permission] !== PermissionStatus.GRANTED) return
    return cloneElement(<Box/>, props)
}