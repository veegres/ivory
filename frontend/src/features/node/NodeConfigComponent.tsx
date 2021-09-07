import {useAppDispatch, useAppSelector} from "../../app/hooks";
import {useEffect} from "react";
import {incrementAsync, selectNodeConfigData} from "./NodeConfigSlice";
import ReactJson from "react-json-view";
import {Grid, Paper} from "@material-ui/core";

export function NodeConfigComponent() {
    const dispatch = useAppDispatch();
    const node = `p4-fr-ppl-1`;
    useEffect(() => { dispatch(incrementAsync(node)) }, [dispatch, node])

    const nodeConfig = useAppSelector(selectNodeConfigData)
    if (!nodeConfig) return null;

    return (
        <Paper elevation={3} style={{ padding: "30px" }}>
            <Grid container spacing={2} direction="column">
                <Grid item style={{fontSize: "20px"}}>Config</Grid>
                <Grid item><ReactJson src={nodeConfig} /></Grid>
            </Grid>
        </Paper>
     )
}
