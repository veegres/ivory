import {Cached} from "@mui/icons-material"
import {Box, Button, Divider, ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material"
import {useState} from "react"

import {Cluster, Node} from "../../../features/cluster/api/type"
import {Feature} from "../../../features/feature"
import {ManageAccess, ManageAccessBox} from "../../../features/management/component/ManageAccess"
import {useRouterQueryList} from "../../../features/query/api/hook"
import {DbConfig, DbPlugin} from "../../../features/query/api/type"
import {Type as QueryType} from "../../../features/query/api/type"
import {Query} from "../../../features/query/component/Query"
import {ErrorDbMissing} from "../../../shared/component/box/ErrorManual"
import {ErrorSmart} from "../../../shared/component/box/ErrorSmart"
import {LinearProgressStateful} from "../../../shared/component/progress/LinearProgressStateful"
import {SxPropsMap} from "../../../shared/helper/type"
import {getQueryConnection} from "../../../shared/helper/utils"
import {useRouterPgCompactTableList} from "../api/hook"
import {PgCompactTableTarget} from "../api/type"
import {PgCompactTableJob} from "./PgCompactTableJob"
import {PgCompactTableJobForm} from "./PgCompactTableJobForm"

const SX: SxPropsMap = {
    loader: {margin: "15px 0"},
    toggle: {display: "flex", flexDirection: "column", alignItems: "center", gap: 1},
    option: {display: "flex", padding: "0 15px", gap: 3},
    form: {flexGrow: 1, display: "flex", flexDirection: "column", justifyContent: "center"},
    refresh: {width: "100%"},
}

enum ListBlock {JOB, QUERY}

type Props = {
    node: Node,
    cluster: Cluster,
}

export function PgCompactTable(props: Props) {
    const {cluster, node} = props
    const [tab, setTab] = useState(ListBlock.JOB)
    const [target, setTarget] = useState<PgCompactTableTarget>()

    const query = useRouterQueryList(QueryType.BLOAT, tab === ListBlock.QUERY)
    const jobs = useRouterPgCompactTableList(cluster.name, tab === ListBlock.JOB)
    const loading = jobs.isFetching || query.isFetching

    if (!node.config.dbPort) return <ErrorDbMissing/>
    const db: DbConfig = {
        plugin: DbPlugin.POSTGRES,
        host: node.config.host,
        port: node.config.dbPort,
        name: target?.database,
        schema: target?.schema,
    }

    return (
        <Box>
            <ManageAccessBox sx={SX.option} feature={Feature.ManageToolPgCompactTableJob} error={true}>
                <Box sx={SX.form}>
                    <PgCompactTableJobForm
                        node={node}
                        cluster={cluster}
                        onClick={() => setTab(ListBlock.JOB)}
                        target={target}
                        setTarget={setTarget}
                    />
                </Box>
                <Divider orientation={"vertical"} flexItem/>
                {renderToggle()}
            </ManageAccessBox>
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
                    <PgCompactTableJob list={jobs.data} cluster={cluster.name} refetchList={jobs.refetch}/>
                )
            case ListBlock.QUERY:
                const queryCon = getQueryConnection(cluster, db.host, db.port)
                return (
                    <ManageAccess feature={Feature.ViewQueryCrudList} error={true}>
                        <Query type={QueryType.BLOAT} connection={{...queryCon, db}}/>
                    </ManageAccess>
                )
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
