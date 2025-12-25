import {Cached} from "@mui/icons-material"
import {Box, Button, Divider, ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material"
import {useState} from "react"

import {useRouterBloatList} from "../../../../api/bloat/hook"
import {BloatTarget} from "../../../../api/bloat/type"
import {Cluster, Instance} from "../../../../api/cluster/type"
import {Permission} from "../../../../api/permission/type"
import {Database, QueryType} from "../../../../api/postgres"
import {useRouterQueryList} from "../../../../api/query/hook"
import {SxPropsMap} from "../../../../app/type"
import {getConnectionRequest} from "../../../../app/utils"
import {ErrorSmart} from "../../../view/box/ErrorSmart"
import {LinearProgressStateful} from "../../../view/progress/LinearProgressStateful"
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
    instance: Instance,
}

export function OverviewBloat(props: Props) {
    const {cluster, instance} = props
    const [tab, setTab] = useState(ListBlock.JOB)
    const [target, setTarget] = useState<BloatTarget>()

    const query = useRouterQueryList(QueryType.BLOAT, tab === ListBlock.QUERY)
    const jobs = useRouterBloatList(cluster.name, tab === ListBlock.JOB)
    const loading = jobs.isFetching || query.isFetching
    const db = {...instance.database, name: target?.database, schema: target?.schema} as Database

    return (
        <Box>
            <AccessBox sx={SX.option} permission={Permission.ManageBloatJob}>
                <Box sx={SX.form}>
                    <OverviewBloatJobForm
                        instance={instance}
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
                return jobs.error ? <ErrorSmart error={jobs.error}/> : <OverviewBloatJob list={jobs.data} cluster={cluster.name} refetchList={jobs.refetch}/>
            case ListBlock.QUERY:
                return <Query type={QueryType.BLOAT} connection={getConnectionRequest(cluster, db)}/>
        }
    }

    function renderToggle() {
        return (
            <Box sx={SX.toggle}>
                <ToggleButtonGroup size={"small"} color={"secondary"} value={tab} orientation={"vertical"}>
                    <ToggleButton value={ListBlock.JOB} onClick={handleJobTab}>
                        Jobs
                    </ToggleButton>
                    <ToggleButton value={ListBlock.QUERY} onClick={handleQueryTab} disabled={!cluster.credentials.postgresId}>
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
