import {PageBox} from "../view/box/PageBox";
import {ErrorSmart} from "../view/box/ErrorSmart";
import {Skeleton, Stack} from "@mui/material";
import {StartupInitialSecret} from "../startup/StartupInitialSecret";
import {StartupSecondary} from "../startup/StartupSecondary";
import {List as ClusterList} from "../cluster/list/List";
import {Overview as ClusterOverview} from "../cluster/overview/Overview";
import {Instance as ClusterInstance} from "../cluster/instance/Instance";
import {AppInfo, SxPropsMap} from "../../type/common";
import {UseQueryResult} from "@tanstack/react-query";
import {StartupLogin} from "../startup/StartupLogin";
import {Menu} from "../settings/menu/Menu";
import {StartupInitialConfig} from "../startup/StartupInitialConfig";

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
    if (!data.secret.ref) return <StartupInitialSecret/>
    if (!data.secret.key) return <StartupSecondary/>
    if (!data.configured) return <StartupInitialConfig/>
    if (!data.auth.authorised) return <StartupLogin type={data.auth.type} error={data.auth.error}/>

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
                <PageBox><ErrorSmart error={error}/></PageBox>
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
            <PageBox>
                <Skeleton variant={"rectangular"} width={"100%"} height={"20vh"} />
            </PageBox>
        )
    }
}
