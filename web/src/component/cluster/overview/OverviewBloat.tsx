import {Box, Button, Divider, ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material";
import {useQuery} from "@tanstack/react-query";
import {bloatApi, queryApi} from "../../../app/api";
import {useState} from "react";
import {CompactTable, QueryType, SxPropsMap} from "../../../app/types";
import {LinearProgressStateful} from "../../view/LinearProgressStateful";
import {TabProps} from "./Overview";
import {OverviewBloatJobForm} from "./OverviewBloatJobForm";
import {OverviewBloatJob} from "./OverviewBloatJob";
import {Query} from "../../shared/query/Query";
import {Cached} from "@mui/icons-material";

const SX: SxPropsMap = {
    loader: {margin: "15px 0"},
    toggle: {display: "flex", flexDirection: "column", alignItems: "center", gap: 1},
    option: {display: "flex", padding: "0 15px", gap: 3},
    form: {flexGrow: 1},
    refresh: {width: "100%"},
}

export function OverviewBloat(props: TabProps) {
    const {cluster, defaultInstance} = props.info
    const [tab, setTab] = useState<"queries" | "jobs">("jobs")
    const [jobs, setJobs] = useState<CompactTable[]>([])

    const query = useQuery(
        ["query", "map", QueryType.BLOAT],
        () => queryApi.map(QueryType.BLOAT),
        {enabled: false},
    )
    const initJobs = useQuery(
        ["instance/bloat/list", cluster.name],
        () => bloatApi.list(cluster.name),
        {onSuccess: (initJobs) => setJobs(initJobs)},
    )
    const loading = initJobs.isFetching || query.isFetching

    return (
        <Box>
            <Box sx={SX.option}>
                <Box sx={SX.form}>
                    <OverviewBloatJobForm
                        defaultInstance={defaultInstance}
                        cluster={cluster}
                        onClick={() => setTab("jobs")}
                        onSuccess={(job) => setJobs([job, ...jobs])}
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
            case "jobs":
                return <OverviewBloatJob list={jobs}/>
            case "queries":
                return <Query type={QueryType.BLOAT} cluster={cluster.name} db={defaultInstance.database}/>
        }
    }

    function renderToggle() {
        return (
            <Box sx={SX.toggle}>
                <ToggleButtonGroup size={"small"} color={"secondary"} value={tab} orientation={"vertical"}>
                    <ToggleButton value={"jobs"} onClick={() => setTab("jobs")}>Jobs</ToggleButton>
                    <ToggleButton value={"queries"} disabled={!cluster.credentials.postgresId} onClick={() => setTab("queries")}>Queries</ToggleButton>
                </ToggleButtonGroup>
                <Tooltip title={`Refetch ${tab}`} placement={"top"} disableInteractive>
                    <Box sx={SX.refresh} component={"span"}>
                        <Button
                            variant={"outlined"}
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

    function handleRefresh() {
        if (tab === "jobs") initJobs.refetch().then()
        else query.refetch().then()
    }
}
