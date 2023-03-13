import {useQuery} from "@tanstack/react-query";
import {queryApi} from "../../../app/api";
import {Database, SxPropsMap} from "../../../type/common";
import {Box} from "@mui/material";
import {ErrorAlert} from "../../view/ErrorAlert";
import React from "react";
import {ChartItem, Color} from "./ChartItem";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", justifyContent: "center", height: "100%", gap: 3},
    row: {display: "flex", flexWrap: "wrap", justifyContent: "space-evenly", gap: 3},
}

type Props = {
    cluster: string,
    db: Database,
}

export function Chart(props: Props) {
    const {cluster: clusterName, db} = props
    const chart = useQuery(
        ["query", "chart", clusterName, db.host, db.port],
        () => queryApi.chart({clusterName, db}),
        {retry: false})
    const {isLoading, data, error} = chart

    if (error) return <ErrorAlert error={error}/>

    return (
        <Box sx={SX.box}>
            <Box sx={SX.row}>
                <ChartItem loading={isLoading} label={"Databases"} value={data?.dbCount} color={Color.DEEP_PURPLE} width={"250px"}/>
                <ChartItem loading={isLoading} label={"Schemas"} value={data?.schemaCount} color={Color.INDIGO} width={"250px"}/>
                <ChartItem loading={isLoading} label={"Connections"} value={data?.connectionCount} color={Color.PURPLE} width={"250px"}/>
                <ChartItem loading={isLoading} label={"Database Uptime"} value={data?.uptime} color={Color.DEEP_ORANGE} width={"250px"}/>
            </Box>
            <Box sx={SX.row}>
                <ChartItem loading={isLoading} label={"Indexes Size"} value={data?.indexSize} color={Color.BLUE_GREY} width={"300px"}/>
                <ChartItem loading={isLoading} label={"Tables Size"} value={data?.tableSize} color={Color.BLUE_GREY} width={"300px"}/>
                <ChartItem loading={isLoading} label={"Database Size"} value={data?.totalSize} color={Color.BLUE_GREY} width={"300px"}/>
            </Box>
        </Box>
    )
}
