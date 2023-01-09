import {Box, Button, Collapse, IconButton, TextField, Tooltip} from "@mui/material";
import {useMutation, useQuery} from "@tanstack/react-query";
import {bloatApi} from "../../../app/api";
import React, {useState} from "react";
import {CompactTable, Style, Target} from "../../../app/types";
import {OverviewBloatJob} from "./OverviewBloatJob";
import {Cached} from "@mui/icons-material";
import {LinearProgressStateful} from "../../view/LinearProgressStateful";
import {TransitionGroup} from "react-transition-group";
import {TabProps} from "./Overview";
import {ClusterNoInstanceError, ClusterNoLeaderError, ClusterNoPostgresPassword} from "./OverviewError";
import {useMutationOptions} from "../../../hook/QueryCustom";

const SX = {
    jobsLoader: {margin: "15px 0"},
    form: {display: "grid", flexGrow: 1, padding: "0 20px", gridTemplateColumns: "repeat(3, 1fr)", gridColumnGap: "30px", gridRowGap: "5px"},
    buttons: {display: "flex", alignItems: "center", gap: 1},
}

const style: Style = {
    transition: {display: "flex", flexDirection: "column", gap: "10px"}
}

export function OverviewBloat({info}: TabProps) {
    const { cluster, instance } = info
    const [target, setTarget] = useState<Target>()
    const [ratio, setRadio] = useState<number>()
    const [jobs, setJobs] = useState<CompactTable[]>([])

    const initJobs = useQuery(['instance/bloat/list', cluster.name], () => bloatApi.list(cluster.name), {
        onSuccess: (initJobs) => setJobs(initJobs)
    })
    const { onError } = useMutationOptions()
    const start = useMutation(bloatApi.start, {
        onSuccess: (job) => setJobs([job, ...jobs]),
        onError,
    })

    if (!instance) return <ClusterNoInstanceError />

    return (
        <Box>
            {renderForm()}
            <LinearProgressStateful sx={SX.jobsLoader} isFetching={initJobs.isFetching || start.isLoading} color={"inherit"} />
            <TransitionGroup style={style.transition}>
                {jobs.map((value) => (
                    <Collapse key={value.uuid}>
                        <OverviewBloatJob key={value.uuid} compactTable={value}/>
                    </Collapse>
                ))}
            </TransitionGroup>
        </Box>
    )

    function renderForm() {
        if (!instance.leader) return <ClusterNoLeaderError />
        if (!cluster.postgresCredId) return <ClusterNoPostgresPassword />

        return (
            <Box sx={SX.form}>
                <TextField
                    size={"small"} label={"Database Name"} variant={"standard"}
                    onChange={(e) => setTarget({...target, dbName: e.target.value})}
                />
                <TextField
                    size={"small"} label="Schema" variant={"standard"}
                    onChange={(e) => setTarget({...target, schema: e.target.value})}
                />
                <TextField
                    size={"small"} label={"Table"} variant={"standard"}
                    onChange={(e) => setTarget({...target, table: e.target.value})}
                />
                <TextField
                    size={"small"} label={"Exclude Schema"} variant={"standard"}
                    onChange={(e) => setTarget({...target, excludeSchema: e.target.value})}
                />
                <TextField
                    size={"small"} label={"Exclude Table"} variant={"standard"}
                    onChange={(e) => setTarget({...target, excludeTable: e.target.value})}
                />
                <Box sx={SX.buttons}>
                    <TextField sx={{ flexGrow: 1 }}
                        size={"small"} label={"Ratio"} type={"number"} variant={"standard"}
                        onChange={(e) => setRadio(parseInt(e.target.value))}
                    />
                    <Box sx={SX.buttons}>
                        <Button variant={"text"} disabled={start.isLoading} onClick={handleRun}>
                            RUN
                        </Button>
                        <Tooltip title={"Refresh list of Jobs"} placement={"top"}>
                            <Box component={"span"}>
                                <IconButton onClick={() => initJobs.refetch()} disabled={initJobs.isFetching}>
                                    <Cached />
                                </IconButton>
                            </Box>
                        </Tooltip>
                    </Box>
                </Box>
            </Box>
        )
    }

    function handleRun() {
        if (instance && cluster?.postgresCredId) {
            const { database: { host, port }} = instance
            start.mutate({
                connection: { host, port, credId: cluster.postgresCredId },
                target,
                ratio,
                cluster: cluster.name
            })
        }
    }
}
