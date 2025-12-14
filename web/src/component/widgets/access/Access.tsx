import {Box, BoxProps} from "@mui/material"
import {cloneElement, Fragment, FragmentProps} from "react"

import {useRouterInfo} from "../../../api/management/hook"
import {Permission, PermissionStatus} from "../../../api/permission/type"

type Props = FragmentProps & {
    permission: Permission,
}

export function Access(props: Props) {
    const info = useRouterInfo()
    const permissions = info.data?.auth.user?.permissions
    if (permissions && permissions[props.permission] !== PermissionStatus.GRANTED) return
    return cloneElement(<Fragment/>, {children: props.children})
}

type PropsBox = BoxProps & {
    permission: Permission,
}

export function AccessBox(props: PropsBox) {
    const info = useRouterInfo()
    const permissions = info.data?.auth.user?.permissions
    if (permissions && permissions[props.permission] !== PermissionStatus.GRANTED) return
    return cloneElement(<Box/>, props)
}