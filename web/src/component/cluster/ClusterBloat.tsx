import {Box, Button, Grid, LinearProgress, TextField} from "@mui/material";
import {useMutation, useQuery} from "react-query";
import {bloatApi} from "../../app/api";
import {useState} from "react";
import {Auth, CompactTable, Target} from "../../app/types";
import {Error} from "../view/Error";
import {ClusterBloatJob} from "./ClusterBloatJob";

const SX = {
    jobsLoader: {minHeight: "4px", margin: "10px 0"},
}

type Props = { leader: string, cluster: string }

export function ClusterBloat({leader, cluster}: Props) {
    const [auth, setAuth] = useState<Auth>({username: '', password: ''})
    const [target, setTarget] = useState<Target>()
    const [ratio, setRadio] = useState<number>()
    const [jobs, setJobs] = useState<CompactTable[]>([])

    const initJobs = useQuery(
        ['node/bloat/list', cluster],
        () => bloatApi.list(cluster),
        {onSuccess: (initJobs) => setJobs(initJobs)}
    )
    const compact = useMutation(bloatApi.start, {onSuccess: (job) => setJobs([job, ...jobs])})

    const [domain, port] = leader.split(":")

    return (
        <Box>
            {domain && port ? renderForm() : <Error error={"No leader found"}/>}
            <Box sx={SX.jobsLoader}>{initJobs.isFetching ? <LinearProgress/> : null}</Box>
            <Grid container gap={2}>
                {jobs.map((value) => <ClusterBloatJob key={value.uuid} compactTable={value}/>)}
            </Grid>
        </Box>
    )

    function renderForm() {
        return (
            <Box>
                <Grid container direction={"row"} justifyContent={"space-between"} alignItems={"center"} flexWrap={"nowrap"}>
                    <Grid item container gap={2}>
                        <TextField
                            required size={"small"} label="Username" variant="standard"
                            onChange={(e) => setAuth({...auth, username: e.target.value})}
                        />
                        <TextField
                            required size={"small"} label="Password" type="password" variant="standard"
                            onChange={(e) => setAuth({...auth, password: e.target.value})}
                        />
                        <TextField
                            size={"small"} label="Ratio" type="number" variant="standard"
                            onChange={(e) => setRadio(parseInt(e.target.value))}
                        />
                    </Grid>
                    <Grid item>
                        <Button variant="contained" size={"small"} disabled={compact.isLoading} onClick={handleRun}>
                            RUN
                        </Button>
                    </Grid>
                </Grid>
                <Grid container gap={2}>
                    <TextField
                        size={"small"} label="Database Name" variant="standard"
                        onChange={(e) => setTarget({...target, dbName: e.target.value})}
                    />
                    <TextField
                        size={"small"} label="Schema" variant="standard"
                        onChange={(e) => setTarget({...target, schema: e.target.value})}
                    />
                    <TextField
                        size={"small"} label="Table" variant="standard"
                        onChange={(e) => setTarget({...target, table: e.target.value})}
                    />
                    <TextField
                        size={"small"} label="Exclude Schema" variant="standard"
                        onChange={(e) => setTarget({...target, excludeSchema: e.target.value})}
                    />
                    <TextField
                        size={"small"} label="Exclude Table" variant="standard"
                        onChange={(e) => setTarget({...target, excludeTable: e.target.value})}
                    />
                </Grid>
            </Box>
        )
    }

    function handleRun() {
        const connection = {host: domain, port: Number(port), ...auth}
        compact.mutate({connection, target, ratio, cluster})
    }
}
