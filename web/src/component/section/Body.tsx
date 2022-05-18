import {Block} from "../view/Block";
import {Error} from "../view/Error";
import {AxiosError} from "axios";
import {CircularProgress, Stack} from "@mui/material";
import {InitialSecret} from "../secret/InitialSecret";
import {RepeatSecret} from "../secret/RepeatSecret";
import {Clusters} from "../clusters/Clusters";
import {Cluster} from "../cluster/Cluster";
import {Node} from "../node/Node";
import React from "react";
import {useQuery} from "react-query";
import {secretApi} from "../../app/api";

export function Body() {
    const { data: status, isLoading, isError, error } = useQuery("secret", secretApi.get)

    if (isError) return <Block><Error error={error as AxiosError}/></Block>
    if (isLoading) return <CircularProgress/>
    if (!status?.ref) return <InitialSecret/>
    if (!status?.key) return <RepeatSecret/>

    return (
        <Stack sx={{ width: "100%", height: "100%" }} spacing={1}>
            <Clusters/>
            <Cluster/>
            <Node/>
        </Stack>
    )
}
