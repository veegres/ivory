import {Box, Grid, Skeleton} from "@mui/material";
import {instanceApi} from "../../../app/api";
import {useQuery} from "@tanstack/react-query";
import {ErrorAlert} from "../../view/ErrorAlert";
import {InstanceColor} from "../../../app/utils";
import {StylePropsMap} from "../../../app/types";
import {useStore} from "../../../provider/StoreProvider";
import {InfoAlert} from "../../view/InfoAlert";
import {PageBlock} from "../../view/PageBlock";

const SX = {
    instanceStatusBlock: {height: "120px", minWidth: "200px", borderRadius: "4px", color: "white", fontSize: "24px", fontWeight: 900},
    item: {margin: "0px 15px"},
    title: {color: "text.secondary", fontWeight: "bold"}
}
const style: StylePropsMap = {
    itemText: {whiteSpace: "pre-wrap"}
}

type ItemProps = { name: string, value?: string }
type StatusProps = { role?: string }

export function Instance() {
    const { store: { activeInstance }, isClusterOverviewOpen } = useStore()
    const {data: instance, isLoading, isError, error} = useQuery(
        ["instance/info", activeInstance?.host, activeInstance?.port],
        () => activeInstance ? instanceApi.info(activeInstance) : undefined,
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
            <Grid container direction={"row"}>
                <Grid item xs={"auto"}>
                    <InstanceStatus role={instance?.role}/>
                </Grid>
                <Grid item xs container direction={"column"}>
                    <Grid item><Item name={"Node"} value={activeInstance?.host}/></Grid>
                    <Grid item><Item name={"Port"} value={activeInstance?.port.toString()}/></Grid>
                    <Grid item><Item name={"State"} value={instance?.state}/></Grid>
                </Grid>
            </Grid>
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
