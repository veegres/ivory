import {Database, QueryType} from "../../../app/types";
import {Skeleton, Stack} from "@mui/material";
import {QueryItem} from "./QueryItem";
import {useQuery} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
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
            {Object.entries(query.data ?? {}).map(([key, value], index) => (
                <QueryItem key={key} add={index === 0} id={key} query={value} cluster={cluster} db={db} type={type}/>
            ))}
        </Stack>
    )

    function renderLoading() {
        return (
            <Stack>
                <Skeleton width={"100%"} height={42}/>
                <Skeleton width={"100%"} height={42}/>
                <Skeleton width={"100%"} height={42}/>
                <Skeleton width={"100%"} height={42}/>
                <Skeleton width={"100%"} height={42}/>
            </Stack>
        )
    }
}
