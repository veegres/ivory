import {UseQueryResult} from "@tanstack/react-query"

import {ClusterBody} from "../../../core/pages/cluster/ClusterBody"
import {ConfigBody} from "../../../core/pages/config/ConfigBody"
import {LoginBody} from "../../../core/pages/login/LoginBody"
import {PermissionsBody} from "../../../core/pages/permission/PermissionsBody"
import {PermissionsLogout} from "../../../core/pages/permission/PermissionsLogout"
import {SecretBodyInitial} from "../../../core/pages/secret/SecretBodyInitial"
import {SecretBodySecondary} from "../../../core/pages/secret/SecretBodySecondary"
import {AppInfo} from "../../../features/management/api/ManagementType"
import {Status} from "../../../features/permission/api/PermissionType"
import {PageErrorBox} from "../box/PageErrorBox"
import {LogoProgress} from "../progress/LogoProgress"

type Props = {
    info: UseQueryResult<AppInfo>,
}

export function Body(props: Props) {
    const {isError, isLoading, data, error} = props.info

    if (isLoading) return <LogoProgress/>
    if (isError) return <PageErrorBox error={error}/>
    if (!data) return <PageErrorBox error={"Something bad happened, we cannot get application initial information"}/>
    if (!data.secret.ref) return <SecretBodyInitial/>
    if (!data.secret.key) return <SecretBodySecondary/>
    if (!data.config.configured || data.config.error) return <ConfigBody configured={data.config.configured} error={data.config.error}/>
    if (!data.auth.authorised) return <LoginBody supported={data.auth.supported} error={data.auth.error}/>
    if (!data.auth.user?.permissions) return <PermissionsLogout username={data.auth.user?.username} error={data.auth.error}/>
    if (!Object.values(data.auth.user.permissions).some(s => s === Status.GRANTED)) {
        const {username, permissions} = data.auth.user
        return <PermissionsBody username={username} permissions={permissions}/>
    }
    return <ClusterBody/>
}
