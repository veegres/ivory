import {Box, BoxProps} from "@mui/material"
import {cloneElement, Fragment, FragmentProps} from "react"

import {ErrorNoAccess} from "../../../shared/component/box/ErrorManual"
import {Feature} from "../../feature"
import {Status} from "../../permission/api/type"
import {useRouterInfo} from "../api/hook"

type Props = FragmentProps & {
    feature: Feature,
    displayError?: boolean,
}

export function ManageAccess(props: Props) {
    const {feature, displayError = false} = props
    const info = useRouterInfo(false)
    const permissions = info.data?.auth.user?.permissions
    if (permissions && permissions[feature] !== Status.GRANTED) return displayError ? <ErrorNoAccess name={feature}/> : undefined
    return cloneElement(<Fragment/>, {children: props.children})
}

type PropsBox = BoxProps & {
    feature: Feature,
    displayError?: boolean,
}

export function ManageAccessBox(props: PropsBox) {
    const {feature, displayError = false} = props
    const info = useRouterInfo(false)
    const permissions = info.data?.auth.user?.permissions
    if (permissions && permissions[feature] !== Status.GRANTED) return displayError ? <ErrorNoAccess name={feature}/> : undefined
    return cloneElement(<Box/>, props)
}