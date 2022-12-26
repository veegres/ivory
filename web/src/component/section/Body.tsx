import {Block} from "../view/Block";
import {ErrorAlert} from "../view/ErrorAlert";
import {Skeleton, Stack} from "@mui/material";
import {InitialSecret} from "../secret/InitialSecret";
import {RepeatSecret} from "../secret/RepeatSecret";
import {Clusters} from "../clusters/Clusters";
import {Cluster} from "../cluster/Cluster";
import {Node} from "../node/Node";
import React from "react";
import {AppInfo} from "../../app/types";
import {UseQueryResult} from "@tanstack/react-query";

const SX = {
    stack: { width: "100%", height: "100%", gap: 1 }
}

type Props = {
    info: UseQueryResult<AppInfo>,
}

export function Body(props: Props) {
    const { isError, isLoading, data, error } = props.info

    if (isLoading) return renderLoading()
    if (isError) return renderError(error)
    if (!data) return renderError("Something bad happened, we cannot get application initial information")
    if (!data.secret.ref) return <InitialSecret/>
    if (!data.secret.key) return <RepeatSecret/>

    return (
        <Stack sx={SX.stack}>
            <Clusters/>
            <Cluster/>
            <Node/>
        </Stack>
    )

    function renderError(error: any) {
        return (
            <Stack sx={SX.stack}>
                <Block><ErrorAlert error={error}/></Block>
            </Stack>
        )
    }

    function renderLoading() {
        return (
            <Stack sx={SX.stack}>
                {renderSkeleton()}
                {renderSkeleton()}
                {renderSkeleton()}
            </Stack>
        )
    }

    function renderSkeleton() {
        return (
            <Block>
                <Skeleton variant={"rectangular"} width={"100%"} height={200} />
            </Block>
        )
    }
}
