import {Button, CircularProgress, Grid, TextField} from "@mui/material";
import {useMutation, useQuery} from "react-query";
import {bloatApi, nodeApi} from "../../app/api";
import {useState} from "react";
import {Auth, CompactTable, Target} from "../../app/types";
import {Error} from "../view/Error";
import {AxiosError} from "axios";
import {NodeJob} from "./NodeJob";

type Props = { node: string }

export function NodeBloat({node}: Props) {
    const [auth, setAuth] = useState<Auth>({username: '', password: ''})
    const [target, setTarget] = useState<Target>()
    const [ratio, setRadio] = useState<number>()
    const [jobs, setJobs] = useState<CompactTable[]>([])

    const initJobs = useQuery(['node/bloat'], bloatApi.list, {onSuccess: (initJobs) => setJobs(initJobs)})
    const leader = useQuery(['node/leader'], () => nodeApi.cluster(node).then((cluster) => {
        return cluster.filter(node => node.role === "leader")[0]
    }))
    const compact = useMutation(bloatApi.add, {onSuccess: (job) => setJobs([job, ...jobs])})

    if (leader.isLoading || initJobs.isLoading) return <CircularProgress/>
    if (leader.isError) return <Error error={leader.error as AxiosError}/>
    if (!leader.data) return <Error error={"No leader found"}/>

    return (
        <Grid container direction={"column"} gap={2}>
            {renderForm()}
            {jobs.map((value) => <NodeJob key={value.uuid} compactTable={value}/>)}
        </Grid>
    )

    function renderForm() {
        return (
            <>
                <Grid item container direction={"row"} justifyContent={"space-between"} alignItems={"center"} flexWrap={"nowrap"}>
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
                        <Button variant="contained" disabled={compact.isLoading} onClick={handleRun}>
                            RUN
                        </Button>
                    </Grid>
                </Grid>
                <Grid item container gap={2}>
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
            </>
        )
    }

    function handleRun() {
        const {data} = leader
        if (data) {
            const connection = {host: data.host, port: data.port, ...auth}
            compact.mutate({connection, target, ratio})
        }
    }
}
