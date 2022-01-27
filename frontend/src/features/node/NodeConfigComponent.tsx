import ReactJson from "react-json-view";
import {Grid, Paper} from "@material-ui/core";
import {getNodeConfig} from "../../app/api";
import { useQuery } from "react-query";

export function NodeConfigComponent() {
    const node = `P4-IO-CHAT-10`;
    const nodeConfig = useQuery(['node/config', node], () => getNodeConfig(node))
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
