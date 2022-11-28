import {Block} from "../view/Block";
import {ErrorAlert} from "../view/ErrorAlert";
import {CircularProgress, Stack} from "@mui/material";
import {InitialSecret} from "../secret/InitialSecret";
import {RepeatSecret} from "../secret/RepeatSecret";
import {Clusters} from "../clusters/Clusters";
import {Cluster} from "../cluster/Cluster";
import {Node} from "../node/Node";
import React from "react";
import {AppInfo} from "../../app/types";
import {UseQueryResult} from "@tanstack/react-query";

type Props = {
    info: UseQueryResult<AppInfo>,
}

export function Body(props: Props) {
    const { isError, isLoading, data, error } = props.info

    if (isError) return <Block><ErrorAlert error={error}/></Block>
    if (isLoading) return <CircularProgress/>
    if (!data) return <Block><ErrorAlert error={"Something bad happened, we cannot get application initial information"}/></Block>
    if (!data.secret.ref) return <InitialSecret/>
    if (!data.secret.key) return <RepeatSecret/>

    return (
        <Stack sx={{ width: "100%", height: "100%", gap: 1 }}>
            <Clusters/>
            <Cluster/>
            <Node/>
        </Stack>
    )
}
