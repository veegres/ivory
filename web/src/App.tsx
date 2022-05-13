import React from 'react';
import {Cluster} from "./component/cluster/Cluster";
import {CircularProgress, Grid, Stack} from "@mui/material";
import {Clusters} from "./component/clusters/Clusters";
import {Header} from "./component/view/Header";
import {Node} from "./component/node/Node"
import {useQuery} from "react-query";
import {secretApi} from "./app/api";
import {InitialSecret} from "./component/secret/InitialSecret";
import {RepeatSecret} from "./component/secret/RepeatSecret";
import {Error} from "./component/view/Error";
import {AxiosError} from "axios";
import {Block} from "./component/view/Block";

export function App() {
    const { data: status, isLoading, isError, error } = useQuery("secret", secretApi.get)

    return (
        <Grid sx={{ height: "100vh" }} container direction={"column"} spacing={2} flexWrap={"nowrap"}>
            <Grid item><Header /></Grid>
            <Grid item container flexGrow={1} justifyContent={"center"} alignItems={"center"}>
                {renderBody()}
            </Grid>
        </Grid>
    );

    function renderBody() {
        if (isError) return <Block><Error error={error as AxiosError} /></Block>
        if (isLoading) return <CircularProgress />
        if (!status?.ref) return <InitialSecret />
        if (!status?.key) return <RepeatSecret />

        return (
            <Stack sx={{ width: "100%", height: "100%" }} spacing={1}>
                <Clusters />
                <Cluster />
                <Node />
            </Stack>
        )
    }
}
