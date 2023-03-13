import {Box, Link, ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material";
import {instanceApi} from "../../../app/api";
import {useQuery} from "@tanstack/react-query";
import {ErrorAlert} from "../../view/ErrorAlert";
import {useStore} from "../../../provider/StoreProvider";
import {InfoAlert} from "../../view/InfoAlert";
import {PageBlock} from "../../view/PageBlock";
import {InstanceInfoStatus} from "./InstanceInfoStatus";
import {InstanceInfoTable} from "./InstanceInfoTable";
import {Database, SxPropsMap} from "../../../type/common";
import {useState} from "react";
import {Chart} from "../../shared/chart/Chart";
import {InstanceTabs} from "../../../type/instance";
import {InstanceMainQueries} from "./InstanceMainQueries";
import {InstanceMainTitle} from "./InstanceMainTitle";
import {ClusterNoPostgresPassword} from "../overview/OverviewError";
import {Add} from "@mui/icons-material";

const SX: SxPropsMap = {
    content: {display: "flex", gap: 3},
    info: {display: "flex", flexDirection: "column", gap: 1},
    main: {flexGrow: 1, overflow: "auto", display: "flex", flexDirection: "column", gap: 2},
    toggle: {padding: "3px"},
}


enum TabsType {QUERY, CHART}
const Tabs: InstanceTabs = {
    [TabsType.CHART]: {
        label: "Simple Charts",
        body: (cluster: string, db: Database) => <Chart cluster={cluster} db={db}/>,
        info: <>
            Here you can check some basic statistic about your database instance.
            If you have some proposal what can be added here, please, suggest
            it <Link href={"https://github.com/veegres/ivory/issues"} target={"_blank"}>here</Link>
        </>
    },
    [TabsType.QUERY]: {
        label: "Queries",
        body: (cluster: string, db: Database, show?: boolean) => <InstanceMainQueries cluster={cluster} db={db} showNew={show!!}/>,
        info: <>
            Here you can run some queries in your postgres to troubleshoot it. Some queries are provided by
            the <i>system</i>, but you can always create your own. As well you're able to
            edit <i>system</i> and <i>custom</i> queries. You can rollback your changes for queries
            at anytime. Be aware that for <i>custom</i> queries the first one is considered as default.
        </>
    }
}


export function Instance() {
    const {store: {activeInstance, activeCluster}, isClusterOverviewOpen} = useStore()
    const [main, setMain] = useState(TabsType.CHART)
    const [showAddQuery, setShowAddQuery] = useState(false)

    const instance = useQuery(
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
        if (instance.isError) return <ErrorAlert error={instance.error}/>
        if (!instance.data && !instance.isLoading) return <ErrorAlert error={"There is no data"}/>

        const {label, info, body} = Tabs[main]
        const {cluster, database} = activeInstance

        const isPasswordExist = activeCluster?.cluster.credentials.postgresId

        return (
            <Box sx={SX.content}>
                <Box sx={SX.info}>
                    <InstanceInfoStatus loading={instance.isLoading} role={instance.data?.role}/>
                    <ToggleButtonGroup size={"small"} color={"secondary"} fullWidth value={main}>
                        <ToggleButton value={TabsType.CHART} onClick={() => setMain(TabsType.CHART)}>Charts</ToggleButton>
                        <ToggleButton value={TabsType.QUERY} onClick={() => setMain(TabsType.QUERY)}>Queries</ToggleButton>
                    </ToggleButtonGroup>
                    <InstanceInfoTable loading={instance.isLoading} instance={instance.data} activeInstance={activeInstance}/>
                </Box>
                <Box sx={SX.main}>
                    <InstanceMainTitle label={label} info={info} renderAdditionalButtons={renderAdditionalButtons()}/>
                    {!isPasswordExist && <ClusterNoPostgresPassword/>}
                    {isPasswordExist && body(cluster, database, showAddQuery)}
                </Box>
            </Box>
        )
    }

    function renderAdditionalButtons() {
        if (main !== TabsType.QUERY) return

        return (
            <ToggleButton
                sx={SX.toggle}
                value={"add"}
                size={"small"}
                selected={showAddQuery}
                onClick={() => setShowAddQuery(!showAddQuery)}
            >
                <Tooltip title={"Add new custom query"} placement={"top"}>
                    <Add/>
                </Tooltip>
            </ToggleButton>
        )
    }
}
