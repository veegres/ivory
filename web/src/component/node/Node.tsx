import {Box, Grid, Skeleton} from "@mui/material";
import {nodeApi} from "../../app/api";
import {useQuery} from "@tanstack/react-query";
import {ErrorAlert} from "../view/ErrorAlert";
import React from "react";
import {InstanceColor} from "../../app/utils";
import {Style} from "../../app/types";
import {useStore} from "../../provider/StoreProvider";
import {InfoAlert} from "../view/InfoAlert";
import {Block} from "../view/Block";

const SX = {
    nodeStatusBlock: {height: '120px', minWidth: '200px', borderRadius: '4px', color: 'white', fontSize: '24px', fontWeight: 900},
    item: {margin: '0px 15px'},
    title: {color: 'text.secondary', fontWeight: 'bold'}
}
const style: Style = {
    itemText: {whiteSpace: 'pre-wrap'}
}

type ItemProps = { name: string, value?: string }
type StatusProps = { role?: string }

export function Node() {
    const { store: { activeInstance }, isClusterOverviewOpen } = useStore()
    const {data: nodePatroni, isLoading, isError, error} = useQuery(
        ['node/overview', activeInstance],
        () => { if (activeInstance) return nodeApi.overview(activeInstance) },
        {enabled: !!activeInstance}
    )

    return (
        <Block withPadding visible={isClusterOverviewOpen()}>
            {renderContent()}
        </Block>
    )

    function renderContent() {
        if (!activeInstance) return <InfoAlert text={"Please, select a node to see the information!"}/>
        if (isError) return <ErrorAlert error={error}/>

        return (
            <Grid container direction="row">
                <Grid item xs="auto">
                    <NodeStatus role={nodePatroni?.role}/>
                </Grid>
                <Grid item xs container direction="column">
                    <Grid item><Item name="Node" value={activeInstance}/></Grid>
                    <Grid item><Item name="State" value={nodePatroni?.state}/></Grid>
                    <Grid item><Item name="Scope" value={nodePatroni?.patroni.scope}/></Grid>
                    <Grid item><Item name="Timeline" value={nodePatroni?.timeline?.toString()}/></Grid>
                    <Grid item><Item name="Xlog" value={JSON.stringify(nodePatroni?.xlog, null, 4)}/></Grid>
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


    function NodeStatus(props: StatusProps) {
        if (isLoading) return <Skeleton variant="rectangular" sx={SX.nodeStatusBlock}/>
        const background = props.role ? InstanceColor[props.role] : undefined

        return (
            <Grid container alignContent="center" justifyContent="center" sx={{...SX.nodeStatusBlock, background}}>
                <Grid item>{nodePatroni?.role.toUpperCase()}</Grid>
            </Grid>
        )
    }
}
