import {Box, BoxProps} from "@mui/material"
import {cloneElement, Fragment, FragmentProps} from "react"

import {Feature} from "../../feature"
import {Status} from "../../permission/api/type"
import {useRouterInfo} from "../api/hook"

type Props = FragmentProps & {
    feature: Feature,
}

export function ManageAccess(props: Props) {
    const info = useRouterInfo(false)
    const permissions = info.data?.auth.user?.permissions
    if (permissions && permissions[props.feature] !== Status.GRANTED) return
    return cloneElement(<Fragment/>, {children: props.children})
}

type PropsBox = BoxProps & {
    feature: Feature,
}

export function ManageAccessBox(props: PropsBox) {
    const info = useRouterInfo(false)
    const permissions = info.data?.auth.user?.permissions
    if (permissions && permissions[props.feature] !== Status.GRANTED) return
    return cloneElement(<Box/>, props)
}