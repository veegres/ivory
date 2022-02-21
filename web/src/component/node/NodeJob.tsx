import {Box, Divider, Grid, IconButton, LinearProgress} from "@mui/material";
import {useEventJob} from "../../app/hooks";
import {useState} from "react";
import {OpenButton} from "../view/OpenButton";
import {CompactTable} from "../../app/types";
import {jobStatus} from "../../app/utils";
import {Clear} from "@mui/icons-material";
import {useMutation, useQueryClient} from "react-query";
import {bloatApi} from "../../app/api";

const SX = {
    console: {fontSize: '13px', width: '100%', background: '#000000D8', padding: '10px 20px', borderRadius: '5px', color: '#d0d0d0'},
    line: {'&:hover': {color: '#f6f6f6'}},
    header: {fontWeight: 'bold', cursor: 'pointer'},
    loader: {margin: '10px 0 5px'},
    divider: {margin: '5px 0'},
    logs: {maxHeight: '350px', overflow: 'auto'},
    button: {padding: '1px', marginLeft: '4px'},
    separator: {marginLeft: '10px'}
}

type Props = { compactTable: CompactTable, logOpen?: boolean }

export function NodeJob({compactTable, logOpen = true}: Props) {
    const {uuid, status: initStatus, command} = compactTable
    const [open, setOpen] = useState(logOpen)
    const {isFetching, logs, status} = useEventJob(uuid, initStatus, open)

    const queryClient = useQueryClient();
    const deleteJob = useMutation(bloatApi.delete, {
        onSuccess: async () => await queryClient.refetchQueries(['node/bloat'])
    })

    const {name, color} = jobStatus[status]
    return (
        <Box sx={SX.console}>
            <Grid container sx={SX.header} onClick={() => setOpen(!open)}>
                <Grid item container justifyContent={"space-between"} flexWrap={"nowrap"}>
                    <Grid item>Command</Grid>
                    <Grid item container xs={"auto"} sx={SX.separator}>
                        <Box sx={{color}}>{name}</Box>
                        <IconButton sx={SX.button} size={"small"} onClick={(e) => {
                            e.stopPropagation();
                            deleteJob.mutate(uuid)
                        }}>
                            <Clear sx={{fontSize: 18}}/>
                        </IconButton>
                    </Grid>
                </Grid>
                <Grid item container justifyContent={"space-between"} flexWrap={"nowrap"}>
                    <Grid item>{command}</Grid>
                    <Grid item container xs={"auto"} sx={SX.separator}>
                        <Box>{uuid}</Box>
                        <Box>
                            <IconButton sx={SX.button} size={"small"}>
                                <OpenButton open={open} size={20}/>
                            </IconButton>
                        </Box>
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
                <Divider sx={SX.divider} textAlign={"left"}>LOGS</Divider>
                <Box display={"p"} sx={SX.logs}>
                    {logs.map((line, index) => (<Box key={index} sx={SX.line}>{line}</Box>))}
                </Box>
                {isFetching ? <LinearProgress sx={SX.loader} color="inherit"/> : null}
            </>
        )
    }
}
