import {Box, BoxProps} from "@mui/material"
import {cloneElement, Fragment, FragmentProps} from "react"

import {Feature} from "../../../features/feature"
import {useRouterInfo} from "../../../features/management/hook"
import {Status} from "../../../features/permission/type"

type Props = FragmentProps & {
    feature: Feature,
}

export function Access(props: Props) {
    const info = useRouterInfo(false)
    const permissions = info.data?.auth.user?.permissions
    if (permissions && permissions[props.feature] !== Status.GRANTED) return
    return cloneElement(<Fragment/>, {children: props.children})
}

type PropsBox = BoxProps & {
    feature: Feature,
}

export function AccessBox(props: PropsBox) {
    const info = useRouterInfo(false)
    const permissions = info.data?.auth.user?.permissions
    if (permissions && permissions[props.feature] !== Status.GRANTED) return
    return cloneElement(<Box/>, props)
}