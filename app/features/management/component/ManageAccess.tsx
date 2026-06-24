import {Box, BoxProps} from "@mui/material"
import {cloneElement, Fragment, FragmentProps} from "react"

import {ErrorNoAccess} from "../../../shared/component/box/ErrorManual"
import {Feature} from "../../feature"
import {Status} from "../../permission/api/type"
import {useRouterInfo} from "../api/hook"

type Props = FragmentProps & {
    feature: Feature,
    error?: boolean,
}

export function ManageAccess(props: Props) {
    const {feature, error = false, ...fragmentProps} = props
    const info = useRouterInfo(false)
    const permissions = info.data?.auth.user?.permissions
    if (permissions && permissions[feature] !== Status.GRANTED) return error ? <ErrorNoAccess name={feature}/> : undefined
    return cloneElement(<Fragment/>, fragmentProps)
}

type PropsBox = BoxProps & {
    feature: Feature,
    error?: boolean,
}

export function ManageAccessBox(props: PropsBox) {
    const {feature, error = false, ...boxProps} = props
    const info = useRouterInfo(false)
    const permissions = info.data?.auth.user?.permissions
    if (permissions && permissions[feature] !== Status.GRANTED) return error ? <ErrorNoAccess name={feature}/> : undefined
    return cloneElement(<Box/>, boxProps)
}