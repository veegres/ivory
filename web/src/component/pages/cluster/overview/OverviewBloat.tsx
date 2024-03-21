import {Box, Button, Divider, ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material";
import {useQuery} from "@tanstack/react-query";
import {BloatApi, QueryApi} from "../../../../app/api";
import {useState} from "react";
import {LinearProgressStateful} from "../../../view/progress/LinearProgressStateful";
import {TabProps} from "./Overview";
import {OverviewBloatJobForm} from "./OverviewBloatJobForm";
import {OverviewBloatJob} from "./OverviewBloatJob";
import {Query} from "../../../shared/query/Query";
import {Cached} from "@mui/icons-material";
import {SxPropsMap} from "../../../../type/common";
import {BloatTarget} from "../../../../type/bloat";
import {QueryType} from "../../../../type/query";

const SX: SxPropsMap = {
    loader: {margin: "15px 0"},
    toggle: {display: "flex", flexDirection: "column", alignItems: "center", gap: 1},
    option: {display: "flex", padding: "0 15px", gap: 3},
    form: {flexGrow: 1, display: "flex", flexDirection: "column", justifyContent: "center"},
    refresh: {width: "100%"},
}

enum ListBlock {JOB, QUERY}

export function OverviewBloat(props: TabProps) {
    const {cluster, defaultInstance} = props.info
    const [tab, setTab] = useState(ListBlock.JOB)
    const [target, setTarget] = useState<BloatTarget>()

    const query = useQuery({
        queryKey: ["query", "map", QueryType.BLOAT],
        queryFn: () => QueryApi.list(QueryType.BLOAT),
        enabled: tab === ListBlock.QUERY,
    })
    const jobs = useQuery({
        initialData: [],
        queryKey: ["instance", "bloat", "list", cluster.name],
        queryFn: () => BloatApi.list(cluster.name),
        enabled: tab === ListBlock.JOB,
    })
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
            <LinearProgressStateful sx={SX.loader} isFetching={loading} color={"inherit"}/>
            {renderBody()}
        </Box>
    )

    function renderBody() {
        switch (tab) {
            case ListBlock.JOB:
                return <OverviewBloatJob list={jobs.data} cluster={cluster.name}/>
            case ListBlock.QUERY:
                return cluster.credentials.postgresId && (
                    <Query type={QueryType.BLOAT} credentialId={cluster.credentials.postgresId} db={db}/>
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
