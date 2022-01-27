import ReactJson from "react-json-view";
import {Grid, Skeleton} from "@mui/material";
import {nodeApi} from "../../app/api";
import { useQuery } from "react-query";
import {useTheme} from "../../provider/ThemeProvider";
import {Error} from "../view/Error";
import React from "react";

export function NodeConfig({ node }: { node: string }) {
    const theme = useTheme();
    const { data: nodeConfig, isLoading, isError, error } = useQuery(['node/config', node], () => nodeApi.config(node))
    if (isError) return <Error error={error} />

    return (
        <Grid container direction="column" style={{ padding: '20px'}}>
            <Grid item style={{fontSize: "20px"}}>Config</Grid>
            <Grid item>
                {isLoading ?
                    <Skeleton variant="rectangular" height={300} /> :
                    <ReactJson
                        src={nodeConfig}
                        collapsed={2}
                        iconStyle="square"
                        displayDataTypes={false}
                        onEdit={() => {}}
                        onDelete={() => {}}
                        onAdd={() => {}}
                        theme={theme.mode === 'dark' ? 'apathy' : 'apathy:inverted'}
                    />
                }
            </Grid>
        </Grid>
     )
}
