import {Database, QueryType} from "../../../app/types";
import {Skeleton, Stack} from "@mui/material";
import {QueryItem} from "./QueryItem";
import {useQuery} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {QueryAdd} from "./QueryAdd";
import {ErrorAlert} from "../../view/ErrorAlert";
import React from "react";

type Props = {
    type: QueryType,
    cluster: string,
    db: Database,
}

export function Query(props: Props) {
    const {type, cluster, db} = props
    const query = useQuery(["query", "map", type], () => queryApi.map(type))

    if (query.isLoading) return renderLoading()
    if (query.error) return <ErrorAlert error={query.error}/>

    return (
        <Stack gap={1}>
            {Object.entries(query.data ?? {}).map(([key, value]) => (
                <QueryItem key={key} id={key} query={value} cluster={cluster} db={db}/>
            ))}
            <QueryAdd type={type}/>
        </Stack>
    )

    function renderLoading() {
        return (
            <>
                <Skeleton width={"100%"} height={42} />
                <Skeleton width={"100%"} height={42} />
                <Skeleton width={"100%"} height={42} />
                <Skeleton width={"100%"} height={42} />
                <Skeleton width={"100%"} height={42} />
            </>
        )
    }
}
