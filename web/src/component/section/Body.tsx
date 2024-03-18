import {PageMainBox} from "../view/box/PageMainBox";
import {Skeleton, Stack} from "@mui/material";
import {SecretBodyInitial} from "../pages/secret/SecretBodyInitial";
import {SecretBodySecondary} from "../pages/secret/SecretBodySecondary";
import {AppInfo, SxPropsMap} from "../../type/common";
import {UseQueryResult} from "@tanstack/react-query";
import {LoginBody} from "../pages/login/LoginBody";
import {ConfigBody} from "../pages/config/ConfigBody";
import {ClusterBody} from "../pages/cluster/ClusterBody";
import {PageErrorBox} from "../view/box/PageErrorBox";

const SX: SxPropsMap = {
    stack: {width: "100%", height: "100%", gap: 4}
}

type Props = {
    info: UseQueryResult<AppInfo>,
}

export function Body(props: Props) {
    const {isError, isLoading, data, error} = props.info

    if (isLoading) return renderLoading()
    if (isError) return <PageErrorBox error={error}/>
    if (!data) return <PageErrorBox error={"Something bad happened, we cannot get application initial information"}/>
    if (!data.secret.ref) return <SecretBodyInitial/>
    if (!data.secret.key) return <SecretBodySecondary/>
    if (!data.configured) return <ConfigBody/>
    if (!data.auth.authorised) return <LoginBody type={data.auth.type} error={data.auth.error}/>
    return <ClusterBody/>

    function renderLoading() {
        return (
            <Stack sx={SX.stack} justifyContent={"space-evenly"}>
                {renderSkeleton()}
                {renderSkeleton()}
                {renderSkeleton()}
            </Stack>
        )
    }

    function renderSkeleton() {
        return (
            <PageMainBox>
                <Skeleton variant={"rectangular"} width={"100%"} height={"20vh"} />
            </PageMainBox>
        )
    }
}
