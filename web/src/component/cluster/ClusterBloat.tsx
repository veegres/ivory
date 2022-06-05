import {Box, Button, Collapse, Grid, IconButton, TextField, Tooltip} from "@mui/material";
import {useMutation, useQuery} from "react-query";
import {bloatApi} from "../../app/api";
import {useState} from "react";
import {Credential, CompactTable, Target, Style} from "../../app/types";
import {Error} from "../view/Error";
import {ClusterBloatJob} from "./ClusterBloatJob";
import {useStore} from "../../provider/StoreProvider";
import {Replay} from "@mui/icons-material";
import {LinearProgressStateful} from "../view/LinearProgressStateful";
import {TransitionGroup} from "react-transition-group";

const SX = {
    jobsLoader: {minHeight: "4px", margin: "10px 0"}
}

const style: Style = {
    transition: {display: "flex", flexDirection: "column", gap: "10px"}
}

export function ClusterBloat() {
    const {store: {activeCluster: {name: cluster, leader}}} = useStore()
    const [auth, setAuth] = useState<Credential>({username: '', password: ''})
    const [target, setTarget] = useState<Target>()
    const [ratio, setRadio] = useState<number>()
    const [jobs, setJobs] = useState<CompactTable[]>([])

    const initJobs = useQuery(
        ['node/bloat/list', cluster],
        () => bloatApi.list(cluster),
        {onSuccess: (initJobs) => setJobs(initJobs)}
    )
    const start = useMutation(bloatApi.start, {onSuccess: (job) => setJobs([job, ...jobs])})

    return (
        <Box>
            {leader ? renderForm() : <Error error={"No leader found"}/>}
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
                            required size={"small"} label={"Username"} variant={"standard"}
                            onChange={(e) => setAuth({...auth, username: e.target.value})}
                        />
                        <TextField
                            required size={"small"} label={"Password"} type={"password"} variant={"standard"}
                            onChange={(e) => setAuth({...auth, password: e.target.value})}
                        />
                        <TextField
                            size={"small"} label={"Ratio"} type={"number"} variant={"standard"}
                            onChange={(e) => setRadio(parseInt(e.target.value))}
                        />
                    </Grid>
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
                    </Grid>
                </Grid>
                <Grid item container width={"auto"} direction={"column"} alignItems={"center"} justifyContent={"space-between"}>
                    <Button variant="contained" size={"small"} disabled={start.isLoading} onClick={handleRun}>
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
        if (leader) {
            const connection = {host: leader.host, port: leader.port, ...auth}
            start.mutate({connection, target, ratio, cluster})
        }
    }
}
