import {Box, Grid, Skeleton} from "@mui/material";
import {instanceApi} from "../../../app/api";
import {useQuery} from "@tanstack/react-query";
import {ErrorAlert} from "../../view/ErrorAlert";
import {InstanceColor} from "../../../app/utils";
import {QueryType, StylePropsMap, SxPropsMap} from "../../../app/types";
import {useStore} from "../../../provider/StoreProvider";
import {InfoAlert} from "../../view/InfoAlert";
import {PageBlock} from "../../view/PageBlock";
import {Query} from "../../shared/query/Query";

const SX: SxPropsMap = {
    instanceStatusBlock: {height: "120px", minWidth: "250px", borderRadius: "4px", color: "white", fontSize: "24px", fontWeight: 900},
    item: {margin: "0px 5px"},
    title: {color: "text.secondary", fontWeight: "bold"},
    content: {display: "flex", gap: 2},
    info: {display: "flex", flexDirection: "column", gap: 1},
    query: {flexGrow: 1, overflow: "auto"},
}
const style: StylePropsMap = {
    itemText: {whiteSpace: "pre-wrap"}
}

type ItemProps = { name: string, value?: string }
type StatusProps = { role?: string }

export function Instance() {
    const {store: {activeInstance}, isClusterOverviewOpen} = useStore()
    const {data: instance, isLoading, isError, error} = useQuery(
        ["instance/info", activeInstance?.sidecar.host, activeInstance?.sidecar.port],
        () => activeInstance ? instanceApi.info({cluster: activeInstance.cluster, ...activeInstance.sidecar}) : undefined,
        {enabled: !!activeInstance}
    )

    return (
        <PageBlock withPadding visible={isClusterOverviewOpen()}>
            {renderContent()}
        </PageBlock>
    )

    function renderContent() {
        if (!activeInstance) return <InfoAlert text={"Please, select a instance to see the information!"}/>
        if (isError) return <ErrorAlert error={error}/>

        return (
            <Box sx={SX.content}>
                <Box sx={SX.info}>
                    <InstanceStatus role={instance?.role}/>
                    <Item name={"State"} value={instance?.state}/>
                    <Item name={"Sidecar"} value={`${activeInstance.sidecar.host}:${activeInstance.sidecar.port.toString()}`}/>
                    <Item name={"Database"} value={`${activeInstance.database.host}:${activeInstance.database.port.toString()}`}/>
                </Box>
                <Box sx={SX.query}>
                    <Query type={QueryType.ACTIVITY} cluster={activeInstance.cluster} db={activeInstance.database}/>
                </Box>
            </Box>
        )
    }

    function Item(props: ItemProps) {
        if (isLoading) return <Skeleton height={24} sx={SX.item}/>

        return (
            <Box sx={SX.item}>
                <Box component={"span"} sx={SX.title}>{props.name}:</Box>
                {` `}
                <span style={style.itemText}>{props.value ?? '-'}</span>
            </Box>
        )
    }


    function InstanceStatus(props: StatusProps) {
        if (isLoading) return <Skeleton variant="rectangular" sx={SX.instanceStatusBlock}/>
        const background = props.role ? InstanceColor[props.role] : undefined

        return (
            <Grid container alignContent="center" justifyContent="center" sx={{...SX.instanceStatusBlock, background}}>
                <Grid item>{instance?.role.toUpperCase()}</Grid>
            </Grid>
        )
    }
}
