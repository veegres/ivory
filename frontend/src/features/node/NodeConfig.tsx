import ReactJson from "react-json-view";
import {Grid} from "@mui/material";
import {nodeApi} from "../../app/api";
import { useQuery } from "react-query";
import {useTheme} from "../../provider/ThemeProvider";

export function NodeConfig() {
    const node = `P4-IO-CHAT-10`;
    const theme = useTheme();
    const { data: nodeConfig } = useQuery(['node/config', node], () => nodeApi.config(node))
    if (!nodeConfig) return null;

    return (
        <Grid container direction="column" style={{ padding: '20px'}}>
            <Grid item style={{fontSize: "20px"}}>Config</Grid>
            <Grid item>
                <ReactJson
                    src={nodeConfig}
                    collapsed={2}
                    iconStyle="square"
                    displayDataTypes={false}
                    onEdit={() => {}}
                    onDelete={() => {}}
                    onAdd={() => {}}
                    theme={theme.mode === 'dark' ? 'apathy' : 'apathy:inverted'}
                />
            </Grid>
        </Grid>
     )
}
