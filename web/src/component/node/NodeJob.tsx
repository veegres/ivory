import {Box, Grid, LinearProgress} from "@mui/material";
import {useEventJob} from "../../app/hooks";
import {useState} from "react";
import {OpenButton} from "../view/OpenButton";

const SX = {
    console: {fontSize: '13px', width: '100%', background: '#000000D8', padding: '10px 20px', borderRadius: '5px', color: '#d0d0d0'},
    line: {'&:hover': {color: '#f6f6f6'}},
    header: {fontWeight: 'bold'},
    arrow: {'&:hover': {background: '#90CAF919'}, margin: '3px 0', cursor: 'pointer'},
}

type Props = { uuid: string, cmd: string }

export function NodeJob({uuid, cmd}: Props) {
    const [open, setOpen] = useState(false)
    const {isFetching, logs} = useEventJob(uuid)

    return (
        <Box sx={SX.console}>
            <Grid container sx={SX.header}>
                <Grid item container justifyContent={"space-between"} flexWrap={"nowrap"} gap={3}>
                    <Grid item>Command</Grid>
                    <Grid item xs={"auto"}>Job ID</Grid>
                </Grid>
                <Grid item container justifyContent={"space-between"} flexWrap={"nowrap"} gap={3}>
                    <Grid item>{cmd}</Grid>
                    <Grid item xs={"auto"}>{uuid}</Grid>
                </Grid>
                <Grid item container sx={SX.arrow} justifyContent={"center"} onClick={() => setOpen(!open)}>
                    <OpenButton open={open} size={20}/>
                </Grid>
            </Grid>
            {renderLogs()}
        </Box>
    )

    function renderLogs() {
        if (!open) return

        return (
            <>
                <Box display={"p"}>
                    {logs.map((line, index) => (<Box key={index} sx={SX.line}>{line}</Box>))}
                </Box>
                {isFetching ? <LinearProgress color="inherit"/> : null}
            </>
        )
    }
}
