import {Box, CircularProgress, Divider, Grid, IconButton, Paper, Tooltip} from "@mui/material";
import {cloneElement, ReactElement, useState} from "react";
import {OpenIcon} from "../../../view/icon/OpenIcon";
import {Clear, Stop} from "@mui/icons-material";
import {shortUuid} from "../../../../app/utils";
import {LinearProgressStateful} from "../../../view/progress/LinearProgressStateful";
import scroll from "../../../../style/scroll.module.css"
import {DynamicRowVirtualizer} from "../../../view/scrolling/DynamicRowVirtualizer";
import {useEventJob} from "../../../../hook/EventJob";
import {SxPropsMap} from "../../../../type/general";
import {Bloat} from "../../../../type/bloat";
import select from "../../../../style/select.module.css";
import {useRouterBloatDelete, useRouterBloatStop} from "../../../../router/bloat";

const SX: SxPropsMap = {
    paper: {fontSize: "13px", width: "100%", padding: "8px 15px"},
    console: {fontSize: "13px", width: "100%", padding: "8px 15px"},
    row: {"&:hover": {color: "primary.main"}},
    emptyLine: {textAlign: "center", textTransform: "uppercase"},
    header: {fontWeight: "bold", cursor: "pointer"},
    loader: {margin: "10px 0 5px"},
    divider: {margin: "5px 0"},
    logs: {colorScheme: "dark"},
    button: {padding: "1px"},
    tooltipBox: {marginLeft: "4px", width: "25px", display: "flex", alignItems: "center", justifyContent: "center"},
    separator: {marginLeft: "10px"},
    credential: {color: "text.disabled", fontWeight: 500}
}

type Props = {
    cluster: string,
    item: Bloat,
}

export function OverviewBloatJobItem(props: Props) {
    const {item, cluster} = props
    const {uuid, status: initStatus, command, credentialId} = item
    const [open, setOpen] = useState(false)
    const {isFetching, logs, status} = useEventJob(uuid, initStatus, open)

    const deleteJob = useRouterBloatDelete(uuid, cluster)
    const stopJob = useRouterBloatStop()

    return (
        <Paper sx={SX.paper} variant={"outlined"}>
            <Grid container sx={SX.header} onClick={() => setOpen(!open)} className={select.none}>
                <Grid item container justifyContent={"space-between"} flexWrap={"nowrap"}>
                    <Grid item container gap={1}>
                        <Box>Command</Box>
                        <Tooltip title={"Postgres Password ID"} placement={"right"}>
                            <Box sx={SX.credential}>{shortUuid(credentialId)}</Box>
                        </Tooltip>
                    </Grid>
                    <Grid item container xs={"auto"} sx={SX.separator}>
                        <Box sx={{color: status.color}}>{status.name}</Box>
                        {status.active ?
                            renderJobButton("Stop", <Stop/>, () => stopJob.mutate(uuid), stopJob.isPending) :
                            renderJobButton("Delete", <Clear/>, () => deleteJob.mutate(uuid), deleteJob.isPending)
                        }
                    </Grid>
                </Grid>
                <Grid item container justifyContent={"space-between"} flexWrap={"nowrap"}>
                    <Grid item>{command}</Grid>
                    <Grid item container xs={"auto"} sx={SX.separator}>
                        <Tooltip title={uuid}><Box>{shortUuid(uuid)}</Box></Tooltip>
                        <Tooltip title={"Open"}>
                            <Box sx={SX.tooltipBox}>
                                <IconButton sx={SX.button} size={"small"}>
                                    <OpenIcon open={open} size={18}/>
                                </IconButton>
                            </Box>
                        </Tooltip>
                    </Grid>
                </Grid>
            </Grid>
            {renderBody()}
        </Paper>
    )

    function renderBody() {
        if (!open) return
        return (
            <>
                <Divider sx={SX.divider} textAlign={"left"}>LOGS</Divider>
                {logs.length === 0 ? isFetching ? (
                    <Box sx={SX.emptyLine}>Waiting for logs</Box>
                ) : (
                    <Box sx={SX.emptyLine}>No logs</Box>
                ) : (
                    <DynamicRowVirtualizer
                        sx={SX.logs}
                        auto={status.active && open}
                        className={scroll.small}
                        sxVirtualRow={SX.row}
                        height={350}
                        rows={logs}
                    />
                )}
                <LinearProgressStateful sx={SX.loader} loading={isFetching} color={"inherit"} line/>
            </>
        )
    }

    function renderJobButton(title: string, icon: ReactElement, onClick: () => void, isLoading: boolean) {
        const fontSize = 18
        return (
            <Tooltip title={title} placement={"top"}>
                <Box sx={SX.tooltipBox}>
                    {isLoading ? <CircularProgress size={fontSize - 3}/> : (
                        <IconButton
                            sx={SX.button}
                            size={"small"}
                            onClick={(e) => {
                                e.stopPropagation();
                                onClick()
                            }}
                        >
                            {cloneElement(icon, {sx: {fontSize}})}
                        </IconButton>
                    )}
                </Box>
            </Tooltip>
        )
    }
}
