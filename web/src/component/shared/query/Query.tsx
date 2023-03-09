import {Database, QueryType, StylePropsMap} from "../../../app/types";
import {Box, Collapse, Skeleton} from "@mui/material";
import {QueryItem} from "./QueryItem";
import {useQuery} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {ErrorAlert} from "../../view/ErrorAlert";
import React from "react";
import {TransitionGroup} from "react-transition-group";

const style: StylePropsMap = {
    transition: {display: "flex", flexDirection: "column", gap: "8px"}
}

type Props = {
    type: QueryType,
    cluster: string,
    db: Database,
}

export function Query(props: Props) {
    const {type, cluster, db} = props
    const query = useQuery(["query", "map", type], () => queryApi.list(type))

    if (query.isLoading) return renderLoading()
    if (query.error) return <ErrorAlert error={query.error}/>

    return (
        <TransitionGroup style={style.transition} appear={false}>
            {(query.data ?? []).map((value) => (
                <Collapse key={value.id}>
                    <QueryItem query={value} cluster={cluster} db={db} type={type}/>
                </Collapse>
            ))}
        </TransitionGroup>
    )

    function renderLoading() {
        return (
            <Box style={style.transition}>
                <Skeleton width={"100%"} height={42}/>
                <Skeleton width={"100%"} height={42}/>
                <Skeleton width={"100%"} height={42}/>
                <Skeleton width={"100%"} height={42}/>
                <Skeleton width={"100%"} height={42}/>
            </Box>
        )
    }
}
