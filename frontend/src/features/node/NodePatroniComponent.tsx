import {Grid} from "@mui/material";
import {nodeApi} from "../../app/api";
import { useQuery } from "react-query";
import {blue, green} from "@mui/material/colors";

export function NodePatroniComponent() {
    const node = `P4-IO-CHAT-13`;
    const { data: nodePatroni } = useQuery(['node/patroni', node], () => nodeApi.patroni(node))
    if (!nodePatroni) return null;

    const background = nodePatroni.role === "master" ? green[500] : blue[500]

    return (
        <Grid container direction="row" sx={{ padding: '10px' }}>
            <Grid item xs="auto">
                {renderNodeStatus(nodePatroni.role)}
            </Grid>
            <Grid item xs container direction="column">
                <Grid item>{renderItem("Node", node)}</Grid>
                <Grid item>{renderItem("State", nodePatroni.state)}</Grid>
                <Grid item>{renderItem("Scope", nodePatroni.patroni.scope)}</Grid>
                <Grid item>{renderItem("Timeline", nodePatroni.timeline.toString())}</Grid>
                <Grid item>{renderItem("Xlog", JSON.stringify(nodePatroni.xlog, null, 4))}</Grid>
            </Grid>
        </Grid>
    )

    function renderItem(name: string, value: string) {
        return (
            <div style={{ padding: '0px 15px' }}>
                <b>{name}:</b>{` `}<span style={{ whiteSpace: 'pre' }}>{value}</span>
            </div>
        )
    }

    function renderNodeStatus(role: string) {
        return (
            <Grid container
                alignContent="center"
                justifyContent="center"
                sx={{ background, color: 'white', fontSize: '24px', height: '120px', minWidth: "200px", fontWeight: 900 }}>
                <Grid item>{role.toUpperCase()}</Grid>
            </Grid>
        )
    }
}
