import {Box, BoxProps} from "@mui/material"
import {cloneElement, Fragment, FragmentProps} from "react"

import {ErrorNoAccess, ErrorNotSupported} from "../../../shared/component/box/ErrorManual"
import {useStore} from "../../../shared/provider/StoreProvider"
import {useRouterClusterOverview} from "../../cluster/api/ClusterHook"
import {Feature} from "../../Feature"
import {PermissionMap, Status} from "../../permission/api/PermissionType"
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
export type Access = "allowed" | "unsupported" | "denied"

type AccessSource = {
    features?: Partial<Record<Feature, boolean>>,
    permissions?: PermissionMap,
}

// useHasAccess is exported for callers that have to tell "unsupported" apart
// from "denied" - a plugin that cannot do a thing at all is worth hiding, while
// a permission the user lacks is worth reporting.
export function useHasAccess(feature: Feature): Access {
    const source = useAccessSource()
    return getAccess(feature, source)
}

// useHasAnyAccess answers whether at least one of the features is allowed. A
// container whose every child hides itself needs it: it cannot tell from the
// children whether it would open onto anything.
export function useHasAnyAccess(features: Feature[]): boolean {
    const source = useAccessSource()
    return features.some(feature => getAccess(feature, source) === "allowed")
}

function useAccessSource(): AccessSource {
    const info = useRouterInfo(false)
    const activeCluster = useStore(s => s.activeCluster)
    const overview = useRouterClusterOverview(activeCluster?.name, false)
    return {features: overview.data?.features, permissions: info.data?.auth.user?.permissions}
}

function getAccess(feature: Feature, source: AccessSource): Access {
    const {features, permissions} = source
    if (features && features[feature] === false) return "unsupported"
    if (permissions && permissions[feature] !== Status.GRANTED) return "denied"
    return "allowed"
}
