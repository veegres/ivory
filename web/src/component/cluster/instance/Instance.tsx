import {Box} from "@mui/material";
import {instanceApi} from "../../../app/api";
import {useQuery} from "@tanstack/react-query";
import {ErrorAlert} from "../../view/ErrorAlert";
import {QueryType, SxPropsMap} from "../../../app/types";
import {useStore} from "../../../provider/StoreProvider";
import {InfoAlert} from "../../view/InfoAlert";
import {PageBlock} from "../../view/PageBlock";
import {Query} from "../../shared/query/Query";
import {InstanceStatus} from "./InstanceStatus";
import {InstanceTable} from "./InstanceTable";

const SX: SxPropsMap = {
    content: {display: "flex", gap: 2},
    info: {display: "flex", flexDirection: "column", gap: 1},
    query: {flexGrow: 1, overflow: "auto"},
}

export function Instance() {
    const {store: {activeInstance}, isClusterOverviewOpen} = useStore()
    const info = useQuery(
        ["instance/info", activeInstance?.sidecar.host, activeInstance?.sidecar.port],
        () => activeInstance ? instanceApi.info({cluster: activeInstance.cluster, ...activeInstance.sidecar}) : undefined,
        {enabled: !!activeInstance}
    )

    return (
        <PageBlock withPadding visible={isClusterOverviewOpen()}>
            {renderContent()}
        </PageBlock>
    )

    function renderContent() {
        if (!activeInstance) return <InfoAlert text={"Please, select a instance to see the information!"}/>
        if (info.isError) return <ErrorAlert error={info.error}/>
        if (!info.data) return <ErrorAlert error={"There is no data"}/>

        return (
            <Box sx={SX.content}>
                <Box sx={SX.info}>
                    <InstanceStatus loading={info.isLoading} role={info.data.role}/>
                    <InstanceTable loading={info.isLoading} instance={info.data} activeInstance={activeInstance}/>
                </Box>
                <Box sx={SX.query}>
                    <Query type={QueryType.ACTIVITY} cluster={activeInstance.cluster} db={activeInstance.database}/>
                </Box>
            </Box>
        )
    }
}
