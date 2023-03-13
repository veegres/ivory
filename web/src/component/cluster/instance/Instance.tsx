import {Box, ToggleButton, ToggleButtonGroup} from "@mui/material";
import {instanceApi} from "../../../app/api";
import {useQuery} from "@tanstack/react-query";
import {ErrorAlert} from "../../view/ErrorAlert";
import {useStore} from "../../../provider/StoreProvider";
import {InfoAlert} from "../../view/InfoAlert";
import {PageBlock} from "../../view/PageBlock";
import {Query} from "../../shared/query/Query";
import {InstanceInfoStatus} from "./InstanceInfoStatus";
import {InstanceInfoTable} from "./InstanceInfoTable";
import {SxPropsMap} from "../../../type/common";
import {QueryType} from "../../../type/query";
import {useState} from "react";
import {Chart} from "../../shared/chart/Chart";
import {ActiveInstance} from "../../../type/instance";

const SX: SxPropsMap = {
    content: {display: "flex", gap: 2},
    info: {display: "flex", flexDirection: "column", gap: 1},
    main: {flexGrow: 1, overflow: "auto"},
}

enum MainBlock {QUERY, CHART}

export function Instance() {
    const {store: {activeInstance}, isClusterOverviewOpen} = useStore()
    const [main, setMain] = useState(MainBlock.CHART)

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
        if (!info.data && !info.isLoading) return <ErrorAlert error={"There is no data"}/>

        return (
            <Box sx={SX.content}>
                <Box sx={SX.info}>
                    <InstanceInfoStatus loading={info.isLoading} role={info.data?.role}/>
                    <ToggleButtonGroup size={"small"} color={"secondary"} fullWidth value={main}>
                        <ToggleButton value={MainBlock.CHART} onClick={() => setMain(MainBlock.CHART)}>Charts</ToggleButton>
                        <ToggleButton value={MainBlock.QUERY} onClick={() => setMain(MainBlock.QUERY)}>Queries</ToggleButton>
                    </ToggleButtonGroup>
                    <InstanceInfoTable loading={info.isLoading} instance={info.data} activeInstance={activeInstance}/>
                </Box>
                <Box sx={SX.main}>
                    {renderMainBlock(activeInstance)}
                </Box>
            </Box>
        )
    }

    function renderMainBlock(instance: ActiveInstance) {
        const {cluster, database} = instance

        switch (main) {
            case MainBlock.CHART:
                return <Chart cluster={cluster} db={database}/>
            case MainBlock.QUERY:
                return <Query type={QueryType.ACTIVITY} cluster={cluster} db={database}/>
        }
    }
}
