import {PageMainBox} from "../view/box/PageMainBox";
import {ErrorSmart} from "../view/box/ErrorSmart";
import {Skeleton, Stack} from "@mui/material";
import {SecretInitial} from "../startup/secret/SecretInitial";
import {SecretSecondary} from "../startup/secret/SecretSecondary";
import {List as ClusterList} from "../cluster/list/List";
import {Overview as ClusterOverview} from "../cluster/overview/Overview";
import {Instance as ClusterInstance} from "../cluster/instance/Instance";
import {AppInfo, SxPropsMap} from "../../type/common";
import {UseQueryResult} from "@tanstack/react-query";
import {Login} from "../startup/login/Login";
import {Menu} from "../settings/menu/Menu";
import {Config} from "../startup/config/Config";

const SX: SxPropsMap = {
    stack: {width: "100%", height: "100%", gap: 4}
}

type Props = {
    info: UseQueryResult<AppInfo>,
}

export function Body(props: Props) {
    const {isError, isLoading, data, error} = props.info

    if (isLoading) return renderLoading()
    if (isError) return renderError(error)
    if (!data) return renderError("Something bad happened, we cannot get application initial information")
    if (!data.secret.ref) return <SecretInitial/>
    if (!data.secret.key) return <SecretSecondary/>
    if (!data.configured) return <Config/>
    if (!data.auth.authorised) return <Login type={data.auth.type} error={data.auth.error}/>

    return (
        <>
            <Menu/>
            <Stack sx={SX.stack}>
                <ClusterList/>
                <ClusterOverview/>
                <ClusterInstance/>
            </Stack>
        </>
    )

    function renderError(error: any) {
        return (
            <Stack sx={SX.stack}>
                <PageMainBox><ErrorSmart error={error}/></PageMainBox>
            </Stack>
        )
    }

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
