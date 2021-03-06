import {Box, CircularProgress, Collapse, Divider, Grid, IconButton, Tooltip} from "@mui/material";
import {useEventJob} from "../../app/hooks";
import React, {ReactElement, useState} from "react";
import {OpenIcon} from "../view/OpenIcon";
import {CompactTable} from "../../app/types";
import {Clear, Stop} from "@mui/icons-material";
import {useMutation, useQueryClient} from "react-query";
import {bloatApi} from "../../app/api";
import {shortUuid} from "../../app/utils";
import {LinearProgressStateful} from "../view/LinearProgressStateful";
import classes from "../../style/scroll.module.css"

const SX = {
    console: {fontSize: "13px", width: "100%", background: "#000", padding: "10px 20px", borderRadius: "5px", color: "#e0e0e0"},
    line: {"&:hover": {color: "#ffffff"}},
    emptyLine: {textAlign: "center"},
    header: {fontWeight: "bold", cursor: "pointer"},
    loader: {margin: "10px 0 5px"},
    divider: {margin: "5px 0"},
    logs: {maxHeight: "350px", overflow: "auto", colorScheme: "dark"},
    button: {padding: "1px", color: '#f6f6f6'},
    tooltipBox: {marginLeft: "4px", width: "25px", display: "flex", alignItems: "center", justifyContent: "center"},
    jobButton: {fontSize: 18},
    separator: {marginLeft: "10px"},
    credential: {color: "#536081", margin: "0 4px"}
}

type Props = { compactTable: CompactTable }

export function ClusterBloatJob({compactTable}: Props) {
    const {uuid, status: initStatus, command, credentialId} = compactTable
    const [open, setOpen] = useState(false)
    const {isFetching, logs, status} = useEventJob(uuid, initStatus, open)

    const queryClient = useQueryClient();
    const deleteJob = useMutation(bloatApi.delete, {
        onSuccess: async () => await queryClient.refetchQueries(['node/bloat/list'])
    })
    const stopJob = useMutation(bloatApi.stop)

    return (
        <Box sx={SX.console}>
            <Grid container sx={SX.header} onClick={() => setOpen(!open)}>
                <Grid item container justifyContent={"space-between"} flexWrap={"nowrap"}>
                    <Grid item container>
                        <Box>Command</Box>
                        <Box sx={SX.credential}>[credential: {shortUuid(credentialId)}]</Box>
                    </Grid>
                    <Grid item container xs={"auto"} sx={SX.separator}>
                        <Box sx={{color: status.color}}>{status.name}</Box>
                        {status.active ?
                            renderJobButton("Stop", <Stop/>, () => stopJob.mutate(uuid), stopJob.isLoading) :
                            renderJobButton("Delete", <Clear/>, () => deleteJob.mutate(uuid), deleteJob.isLoading)
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
            <Collapse in={open}>
                <Divider sx={SX.divider} textAlign={"left"} light>LOGS</Divider>
                <Box display={"p"} sx={SX.logs} className={classes.scroll}>
                    {logs.length === 0 ? isFetching ? (
                        <Box sx={SX.emptyLine}>{"< WAITING FOR LOGS >"}</Box>
                    ) : (
                        <Box sx={SX.emptyLine}>{"< NO LOGS >"}</Box>
                    ) : logs.map((line, index) => (
                        <Box key={index} sx={SX.line}>{line}</Box>
                    ))}
                </Box>
                <LinearProgressStateful sx={SX.loader} isFetching={isFetching} color={"inherit"} line />
            </Collapse>
        </Box>
    )

    function renderJobButton(title: string, icon: ReactElement, onClick: () => void, isLoading: boolean) {
        return (
            <Tooltip title={title} placement={"top"}>
                <Box sx={SX.tooltipBox}>
                    {isLoading ? <CircularProgress size={SX.jobButton.fontSize - 3}/> : (
                        <IconButton
                            sx={SX.button}
                            size={"small"}
                            onClick={(e) => {e.stopPropagation(); onClick()}}
                        >
                            {React.cloneElement(icon, {sx: SX.jobButton})}
                        </IconButton>
                    )}
                </Box>
            </Tooltip>
        )
    }
}
