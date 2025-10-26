import {UseQueryResult} from "@tanstack/react-query"

import {AppInfo} from "../../api/management/type"
import {ClusterBody} from "../../component/pages/cluster/ClusterBody"
import {ConfigBody} from "../../component/pages/config/ConfigBody"
import {LoginBody} from "../../component/pages/login/LoginBody"
import {SecretBodyInitial} from "../../component/pages/secret/SecretBodyInitial"
import {SecretBodySecondary} from "../../component/pages/secret/SecretBodySecondary"
import {PageErrorBox} from "../../component/view/box/PageErrorBox"
import {LogoProgress} from "../../component/view/progress/LogoProgress"

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
    return <ClusterBody/>
}
