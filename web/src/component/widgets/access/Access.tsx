import {Box, BoxProps} from "@mui/material"
import {cloneElement, Fragment, FragmentProps} from "react"

import {Feature} from "../../../api/feature"
import {useRouterInfo} from "../../../api/management/hook"
import {Status} from "../../../api/permission/type"

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