import {Box, Button, Collapse, Grid, IconButton, TextField, Tooltip} from "@mui/material";
import {useMutation, useQuery} from "react-query";
import {bloatApi} from "../../app/api";
import React, {useState} from "react";
import {CompactTable, Style, Target} from "../../app/types";
import {ClusterBloatJob} from "./ClusterBloatJob";
import {Replay} from "@mui/icons-material";
import {LinearProgressStateful} from "../view/LinearProgressStateful";
import {TransitionGroup} from "react-transition-group";
import {TabProps} from "./Cluster";
import {ClusterNoInstanceError, ClusterNoLeaderError} from "./ClusterError";

const SX = {
    jobsLoader: {minHeight: "4px", margin: "10px 0"}
}

const style: Style = {
    transition: {display: "flex", flexDirection: "column", gap: "10px"}
}

export function ClusterBloat({info}: TabProps) {
    const { cluster, instance } = info
    const [target, setTarget] = useState<Target>()
    const [ratio, setRadio] = useState<number>()
    const [jobs, setJobs] = useState<CompactTable[]>([])

    const initJobs = useQuery(
        ['node/bloat/list', cluster.name],
        () => bloatApi.list(cluster.name),
        {onSuccess: (initJobs) => setJobs(initJobs)}
    )
    const start = useMutation(bloatApi.start, {onSuccess: (job) => setJobs([job, ...jobs])})

    if (!instance) return <ClusterNoInstanceError />

    return (
        <Box>
            {instance.leader ? renderForm() : <ClusterNoLeaderError />}
            <LinearProgressStateful sx={SX.jobsLoader} isFetching={initJobs.isFetching || start.isLoading} />
            <TransitionGroup style={style.transition}>
                {jobs.map((value) => (
                    <Collapse key={value.uuid}>
                        <ClusterBloatJob key={value.uuid} compactTable={value}/>
                    </Collapse>
                ))}
            </TransitionGroup>
        </Box>
    )

    function renderForm() {
        return (
            <Grid container justifyContent={"space-between"} flexWrap={"nowrap"}>
                <Grid item container flexGrow={1} direction={"column"} alignItems={"center"} >
                    <Grid item container gap={2}>
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
                        <TextField
                            size={"small"} label={"Ratio"} type={"number"} variant={"standard"}
                            onChange={(e) => setRadio(parseInt(e.target.value))}
                        />
                    </Grid>
                </Grid>
                <Grid item container width={"auto"} direction={"column"} alignItems={"center"} justifyContent={"space-between"}>
                    <Button variant={"contained"} disabled={start.isLoading || !cluster.postgresCredId} onClick={handleRun}>
                        RUN
                    </Button>
                    <Tooltip title={"Reload Jobs"} placement={"left"}>
                        <Box component={"span"}>
                            <IconButton onClick={() => initJobs.refetch()} disabled={initJobs.isFetching}>
                                <Replay/>
                            </IconButton>
                        </Box>
                    </Tooltip>
                </Grid>
            </Grid>
        )
    }

    function handleRun() {
        if (instance && cluster?.postgresCredId) {
            const { host, port } = instance
            start.mutate({
                connection: { host, port, credId: cluster.postgresCredId },
                target,
                ratio,
                cluster: cluster.name
            })
        }
    }
}
