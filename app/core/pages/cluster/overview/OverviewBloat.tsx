import {Cached} from "@mui/icons-material"
import {Box, Button, Divider, ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material"
import {useState} from "react"

import {useRouterBloatList} from "../../../../features/bloat/hook"
import {BloatTarget} from "../../../../features/bloat/type"
import {Cluster, Node} from "../../../../features/cluster/type"
import {Config, Plugin as DbPlugin} from "../../../../features/database/type"
import {Feature} from "../../../../features/feature"
import {useRouterQueryList} from "../../../../features/query/hook"
import {Type as QueryType} from "../../../../features/query/type"
import {ErrorDbMissing} from "../../../../shared/component/box/ErrorManual"
import {ErrorSmart} from "../../../../shared/component/box/ErrorSmart"
import {LinearProgressStateful} from "../../../../shared/component/progress/LinearProgressStateful"
import {SxPropsMap} from "../../../../shared/helper/type"
import {getQueryConnection} from "../../../../shared/helper/utils"
import {AccessBox} from "../../../widgets/access/Access"
import {Query} from "../../../widgets/query/Query"
import {OverviewBloatJob} from "./OverviewBloatJob"
import {OverviewBloatJobForm} from "./OverviewBloatJobForm"

const SX: SxPropsMap = {
    loader: {margin: "15px 0"},
    toggle: {display: "flex", flexDirection: "column", alignItems: "center", gap: 1},
    option: {display: "flex", padding: "0 15px", gap: 3},
    form: {flexGrow: 1, display: "flex", flexDirection: "column", justifyContent: "center"},
    refresh: {width: "100%"},
}

enum ListBlock {JOB, QUERY}

type Props = {
    cluster: Cluster,
    node: Node,
}

export function OverviewBloat(props: Props) {
    const {cluster, node} = props
    const [tab, setTab] = useState(ListBlock.JOB)
    const [target, setTarget] = useState<BloatTarget>()

    const query = useRouterQueryList(QueryType.BLOAT, tab === ListBlock.QUERY)
    const jobs = useRouterBloatList(cluster.name, tab === ListBlock.JOB)
    const loading = jobs.isFetching || query.isFetching

    if (!node.config.dbPort) return <ErrorDbMissing/>
    const db: Config = {
        plugin: DbPlugin.POSTGRES,
        host: node.config.host,
        port: node.config.dbPort,
        name: target?.database,
        schema: target?.schema,
    }

    return (
        <Box>
            <AccessBox sx={SX.option} feature={Feature.ManageToolBloatJob}>
                <Box sx={SX.form}>
                    <OverviewBloatJobForm
                        node={node}
                        cluster={cluster}
                        onClick={() => setTab(ListBlock.JOB)}
                        target={target}
                        setTarget={setTarget}
                    />
                </Box>
                <Divider orientation={"vertical"} flexItem/>
                {renderToggle()}
            </AccessBox>
            <LinearProgressStateful sx={SX.loader} loading={loading} color={"inherit"}/>
            {renderBody()}
        </Box>
    )

    function renderBody() {
        switch (tab) {
            case ListBlock.JOB:
                return jobs.error ? (
                    <ErrorSmart error={jobs.error}/>
                ) : (
                    <OverviewBloatJob list={jobs.data} cluster={cluster.name} refetchList={jobs.refetch}/>
                )
            case ListBlock.QUERY:
                const queryCon = getQueryConnection(cluster, db.host, db.port)
                return <Query type={QueryType.BLOAT} connection={{...queryCon, db}}/>
        }
    }

    function renderToggle() {
        return (
            <Box sx={SX.toggle}>
                <ToggleButtonGroup size={"small"} color={"secondary"} value={tab} orientation={"vertical"}>
                    <ToggleButton value={ListBlock.JOB} onClick={handleJobTab}>
                        Jobs
                    </ToggleButton>
                    <ToggleButton value={ListBlock.QUERY} onClick={handleQueryTab} disabled={!cluster.vaults.databaseId}>
                        Queries
                    </ToggleButton>
                </ToggleButtonGroup>
                <Tooltip title={`Refetch ${ListBlock[tab]}`} placement={"top"} disableInteractive>
                    <Box sx={SX.refresh} component={"span"}>
                        <Button
                            color={"secondary"}
                            fullWidth
                            size={"small"}
                            disabled={loading}
                            onClick={handleRefresh}
                        >
                            <Cached/>
                        </Button>
                    </Box>
                </Tooltip>
            </Box>
        )
    }

    function handleJobTab() {
        setTab(ListBlock.JOB)
        jobs.refetch().then()
    }

    function handleQueryTab() {
        setTab(ListBlock.QUERY)
        query.refetch().then()
    }

    function handleRefresh() {
        if (tab === ListBlock.JOB) jobs.refetch().then()
        else query.refetch().then()
    }
}
