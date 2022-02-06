import {Box, Button, Grid, LinearProgress, TextField} from "@mui/material";
import {useQuery} from "react-query";
import {cliApi, nodeApi} from "../../app/api";
import {useState} from "react";
import {Connection, Target} from "../../app/types";

const SX = {
    console: {width: '100%', background: '#000000D8', padding: '20px', borderRadius: '5px', color: '#d0d0d0'},
    line: {'&:hover': {color: 'white'}}
}

type Props = { node: string }

export function NodeBloat({node}: Props) {
    const [connection, setConnection] = useState<Connection>({host: '', port: 0, username: '', password: ''})
    const [target, setTarget] = useState<Target>()
    const [ratio, setRadio] = useState<number>()

    useQuery(['node/cluster'], () => nodeApi.cluster(node), {
        onSuccess: data => data.forEach(node => {
            if (node.role === "leader") setConnection({...connection, host: node.host, port: node.port})
        })
    })
    const compact = useQuery(
        ['cli/pgcompacttable', node],
        () => cliApi.pgcompacttable({connection, target, ratio}),
        {enabled: false}
    )

    return (
        <Grid container direction={"column"} gap={2}>
            {renderForm()}
            {renderConsole(compact.isFetching, compact.data)}
        </Grid>
    )

    function renderForm() {
        return (
            <>
                <Grid item container direction={"row"} justifyContent={"space-between"} alignItems={"center"}>
                    <Grid item display={"flex"} gap={2}>
                        <TextField
                            required size={"small"} label="Username" variant="standard"
                            onChange={(e) => setConnection({...connection, username: e.target.value})}
                        />
                        <TextField
                            required size={"small"} label="Password" type="password" variant="standard"
                            onChange={(e) => setConnection({...connection, password: e.target.value})}
                        />
                        <TextField
                            size={"small"} label="Ratio" type="number" variant="standard"
                            onChange={(e) => setRadio(parseInt(e.target.value))}
                        />
                    </Grid>
                    <Grid item>
                        <Button variant="contained" size={"small"} disabled={compact.isFetching} onClick={() => compact.refetch()}>
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

    function renderConsole(isFetching: boolean, lines?: string[]) {
        if (!lines && !isFetching) return null

        return (
            <Box display={"p"} sx={SX.console}>
                {isFetching ? <LinearProgress color="inherit"/> : (
                    lines!!.map((line, index) => (<Box key={index} sx={SX.line}>{line}</Box>))
                )}
            </Box>
        )
    }
}
