import {SecretBodyInitial} from "../../component/pages/secret/SecretBodyInitial";
import {SecretBodySecondary} from "../../component/pages/secret/SecretBodySecondary";
import {AppInfo} from "../../api/management/type";
import {UseQueryResult} from "@tanstack/react-query";
import {LoginBody} from "../../component/pages/login/LoginBody";
import {ConfigBody} from "../../component/pages/config/ConfigBody";
import {ClusterBody} from "../../component/pages/cluster/ClusterBody";
import {PageErrorBox} from "../../component/view/box/PageErrorBox";
import {LogoProgress} from "../../component/view/progress/LogoProgress";

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
    if (!data.config.configured) return <ConfigBody error={data.config.error}/>
    if (!data.auth.authorised) return <LoginBody type={data.auth.type} error={data.auth.error}/>
    return <ClusterBody/>
}
