import {UseQueryResult} from "@tanstack/react-query"

import {ClusterBody} from "../../../core/pages/cluster/ClusterBody"
import {ConfigBody} from "../../../core/pages/config/ConfigBody"
import {LoginBody} from "../../../core/pages/login/LoginBody"
import {PermissionsBody} from "../../../core/pages/permission/PermissionsBody"
import {PermissionsLogout} from "../../../core/pages/permission/PermissionsLogout"
import {SecretBodyInitial} from "../../../core/pages/secret/SecretBodyInitial"
import {SecretBodySecondary} from "../../../core/pages/secret/SecretBodySecondary"
import {UserRegistrationBody} from "../../../core/pages/user/UserRegistrationBody"
import {AppInfo} from "../../../features/management/api/ManagementType"
import {Status} from "../../../features/permission/api/PermissionType"
import {getRouteTail, isHomePath, redirectToHome} from "../../helper/HelperUrl"
import {PageErrorBox} from "../box/PageErrorBox"
import {LogoProgress} from "../progress/LogoProgress"

type Props = {
    info: UseQueryResult<AppInfo>,
}

export function Body(props: Props) {
    const {isError, isLoading, data, error} = props.info

    // 1. Check all errors
    if (isLoading) return <LogoProgress/>
    if (isError) return <PageErrorBox error={error}/>
    if (!data) return <PageErrorBox error={"Something bad happened, we cannot get application initial information"}/>

    // 2. Check if we need to configure Ivory
    if (!data.secret.ref) return <SecretBodyInitial/>
    if (!data.secret.key) return <SecretBodySecondary/>
    if (!data.config.configured || data.config.error) return <ConfigBody configured={data.config.configured} error={data.config.error}/>

    // 3. Check if user want to register (it has its own logic after it)
    const registration = getRouteTail("/user")
    if (registration) return renderUserRegistration(registration.tail)

    // 4. Check user state
    if (!data.auth.authorised) return <LoginBody supported={data.auth.supported} error={data.auth.error}/>
    if (!data.auth.user?.permissions) return <PermissionsLogout username={data.auth.user?.username} error={data.auth.error}/>
    if (!Object.values(data.auth.user.permissions).some(s => s === Status.GRANTED)) {
        const {username, permissions} = data.auth.user
        return <PermissionsBody username={username} permissions={permissions}/>
    }

    // 5. Check route, if not root go to root
    if (!isHomePath()) return renderRoutes()

    // 6. Render main body
    return <ClusterBody/>


    function renderRoutes() {
        return renderHome()
    }

    function renderUserRegistration(token: string) {
        if (!token) return renderHome()
        return <UserRegistrationBody token={token}/>
    }

    function renderHome() {
        redirectToHome()
        return <LogoProgress/>
    }

}
