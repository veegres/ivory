import {Box, BoxProps} from "@mui/material"
import {cloneElement, Fragment, FragmentProps} from "react"

import {ErrorNoAccess, ErrorNotSupported} from "../../../shared/component/box/ErrorManual"
import {useStore} from "../../../shared/provider/StoreProvider"
import {useRouterClusterOverview} from "../../cluster/api/ClusterHook"
import {Feature} from "../../Feature"
import {Status} from "../../permission/api/PermissionType"
import {useRouterInfo} from "../api/ManagementHook"

type Props = FragmentProps & {
    feature: Feature,
    error?: boolean,
}

export function ManageAccess(props: Props) {
    const {feature, error = false, ...fragmentProps} = props
    const access = useHasAccess(feature)
    if (access !== "allowed") return error ? renderError(feature, access) : undefined
    return cloneElement(<Fragment/>, fragmentProps)
}

type PropsBox = BoxProps & {
    feature: Feature,
    error?: boolean,
}

export function ManageAccessBox(props: PropsBox) {
    const {feature, error = false, ...boxProps} = props
    const access = useHasAccess(feature)
    if (access !== "allowed") return error ? renderError(feature, access) : undefined
    return cloneElement(<Box/>, boxProps)
}

function renderError(feature: Feature, access: Access) {
    return access === "unsupported" ? <ErrorNotSupported name={feature}/> : <ErrorNoAccess name={feature}/>
}

// NOTE: "unsupported" is returned without ever looking at permissions, since an
// unsupported feature stays denied no matter the permission outcome would have been
type Access = "allowed" | "unsupported" | "denied"

function useHasAccess(feature: Feature): Access {
    const info = useRouterInfo(false)
    const activeCluster = useStore(s => s.activeCluster)
    const overview = useRouterClusterOverview(activeCluster?.name, false)

    const clusterFeatures = overview.data?.features
    if (clusterFeatures && clusterFeatures[feature] === false) return "unsupported"

    const permissions = info.data?.auth.user?.permissions
    if (permissions && permissions[feature] !== Status.GRANTED) return "denied"

    return "allowed"
}
