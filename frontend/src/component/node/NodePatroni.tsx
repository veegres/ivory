import {Box, Grid, Skeleton} from "@mui/material";
import {nodeApi} from "../../app/api";
import { useQuery } from "react-query";
import {blue, green} from "@mui/material/colors";
import {Error} from "../view/Error";
import React from "react";
import {AxiosError} from "axios";

const SX = {
    nodeStatusBlock: { height: '120px', minWidth: '200px', borderRadius: '4px' },
    item: { margin: '0px 15px' }
}
export function NodePatroni({ node }: { node: string }) {
    const { data: nodePatroni, isLoading, isError, error } = useQuery(['node/patroni', node], () => nodeApi.patroni(node))
    if (isError) return <Error error={error as AxiosError} />

    return (
        <Grid container direction="row">
            <Grid item xs="auto">
                <NodeStatus />
            </Grid>
            <Grid item xs container direction="column">
                <Grid item><Item name="Node" value={node} /></Grid>
                <Grid item><Item name="State" value={nodePatroni?.state} /></Grid>
                <Grid item><Item name="Scope" value={nodePatroni?.patroni.scope} /></Grid>
                <Grid item><Item name="Timeline" value={nodePatroni?.timeline.toString()} /></Grid>
                <Grid item><Item name="Xlog" value={JSON.stringify(nodePatroni?.xlog, null, 4)} /></Grid>
            </Grid>
        </Grid>
    )

    function Item(props: { name: string, value?: string }) {
        if (isLoading) return <Skeleton height={24} sx={SX.item} />

        return (
            <Box sx={SX.item}>
                <b>{props.name}:</b>{` `}<span style={{ whiteSpace: 'pre-wrap' }}>{props.value ?? '-'}</span>
            </Box>
        )
    }

    function NodeStatus() {
        if (isLoading) return <Skeleton variant="rectangular" sx={SX.nodeStatusBlock} />
        const background = nodePatroni?.role === "master" ? green[500] : blue[500]

        return (
            <Grid container
                alignContent="center"
                justifyContent="center"
                sx={{ ...SX.nodeStatusBlock, background, color: 'white', fontSize: '24px', fontWeight: 900 }}>
                <Grid item>{nodePatroni?.role.toUpperCase()}</Grid>
            </Grid>
        )
    }
}
