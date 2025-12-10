import {Cached} from "@mui/icons-material"
import {Box, Button, Divider, ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material"
import {useState} from "react"

import {useRouterBloatList} from "../../../../api/bloat/hook"
import {BloatTarget} from "../../../../api/bloat/type"
import {ActiveCluster} from "../../../../api/cluster/type"
import {useRouterQueryList} from "../../../../api/query/hook"
import {QueryType} from "../../../../api/query/type"
import {SxPropsMap} from "../../../../app/type"
import {getQueryConnection} from "../../../../app/utils"
import {ErrorSmart} from "../../../view/box/ErrorSmart"
import {LinearProgressStateful} from "../../../view/progress/LinearProgressStateful"
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
    info: ActiveCluster
}

export function OverviewBloat(props: Props) {
    const {cluster, defaultInstance} = props.info
    const [tab, setTab] = useState(ListBlock.JOB)
    const [target, setTarget] = useState<BloatTarget>()

    const query = useRouterQueryList(QueryType.BLOAT, tab === ListBlock.QUERY)
    const jobs = useRouterBloatList(cluster.name, tab === ListBlock.JOB)
    const loading = jobs.isFetching || query.isFetching
    const db = {...defaultInstance.database, database: target?.dbName}

    return (
        <Box>
            <Box sx={SX.option}>
                <Box sx={SX.form}>
                    <OverviewBloatJobForm
                        defaultInstance={defaultInstance}
                        cluster={cluster}
                        onClick={() => setTab(ListBlock.JOB)}
                        target={target}
                        setTarget={setTarget}
                    />
                </Box>
                <Divider orientation={"vertical"} flexItem/>
                {renderToggle()}
            </Box>
            <LinearProgressStateful sx={SX.loader} loading={loading} color={"inherit"}/>
            {renderBody()}
        </Box>
    )

    function renderBody() {
        switch (tab) {
            case ListBlock.JOB:
                return jobs.error ? <ErrorSmart error={jobs.error}/> : <OverviewBloatJob list={jobs.data} cluster={cluster.name}/>
            case ListBlock.QUERY:
                return <Query type={QueryType.BLOAT} connection={getQueryConnection(cluster, db)}/>
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
