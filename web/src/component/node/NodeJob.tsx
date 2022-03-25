import {Box, CircularProgress, Divider, Grid, IconButton, LinearProgress, Tooltip} from "@mui/material";
import {useEventJob} from "../../app/hooks";
import React, {ReactElement, useState} from "react";
import {OpenIcon} from "../view/OpenIcon";
import {CompactTable} from "../../app/types";
import {Clear, Stop} from "@mui/icons-material";
import {useMutation, useQueryClient} from "react-query";
import {bloatApi} from "../../app/api";

const SX = {
    console: {fontSize: '13px', width: '100%', background: '#000000D8', padding: '10px 20px', borderRadius: '5px', color: '#e0e0e0'},
    line: {'&:hover': {color: '#ffffff'}},
    emptyLine: {textAlign: 'center'},
    header: {fontWeight: 'bold', cursor: 'pointer'},
    loader: {margin: '10px 0 5px'},
    divider: {margin: '5px 0'},
    logs: {maxHeight: '350px', overflow: 'auto'},
    button: {padding: '1px', color: '#f6f6f6'},
    tooltipBox: {marginLeft: '4px'},
    jobButton: {fontSize: 18},
    separator: {marginLeft: '10px'}
}

type Props = { compactTable: CompactTable }

export function NodeJob({compactTable}: Props) {
    const {uuid, status: initStatus, command} = compactTable
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
                    <Grid item>Command</Grid>
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
                        <Tooltip title={uuid}><Box>{uuid.substring(0, 8)}</Box></Tooltip>
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
            {renderLogs()}
        </Box>
    )

    function renderLogs() {
        if (!open) return

        return (
            <>
                <Divider sx={SX.divider} textAlign={"left"} light={true}>LOGS</Divider>
                <Box display={"p"} sx={SX.logs}>
                    {logs.length === 0 ? isFetching ? (
                        <Box sx={SX.emptyLine}>{"< WAITING FOR LOGS >"}</Box>
                    ) : (
                        <Box sx={SX.emptyLine}>{"< NO LOGS >"}</Box>
                    ) : logs.map((line, index) => (
                        <Box key={index} sx={SX.line}>{line}</Box>
                    ))}
                </Box>
                {isFetching ? <LinearProgress sx={SX.loader} color="inherit"/> : null}
            </>
        )
    }

    function renderJobButton(title: string, icon: ReactElement, onClick: () => void, isLoading: boolean) {
        return (
            <Tooltip title={title} placement={"top"}>
                <Box sx={SX.tooltipBox}>
                    {isLoading ? <CircularProgress size={SX.jobButton.fontSize - 2} /> : (
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
