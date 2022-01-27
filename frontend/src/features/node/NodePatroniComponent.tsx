import {Grid} from "@mui/material";
import css from "./NodePatroni.module.css"
import classNames from "classnames";
import {nodeApi} from "../../app/api";
import { useQuery } from "react-query";

export function NodePatroniComponent() {
    const node = `P4-IO-CHAT-10`;
    const { data: nodePatroni } = useQuery(['node/patroni', node], () => nodeApi.patroni(node))
    if (!nodePatroni) return null;

    const roleClassName = classNames(css.role, { [css.master]: nodePatroni.role === "master" })

    return (
        <Grid container direction="row">
            <Grid className={roleClassName} item xs={2} container alignContent="center" justifyContent="center">
                <Grid item>{nodePatroni.role.toUpperCase()}</Grid>
            </Grid>
            <Grid item xs={10} container direction="column">
                <Grid item>{renderItem("Node", node)}</Grid>
                <Grid item>{renderItem("State", nodePatroni.state)}</Grid>
                <Grid item>{renderItem("Scope", nodePatroni.patroni.scope)}</Grid>
                <Grid item>{renderItem("Timeline", nodePatroni.timeline.toString())}</Grid>
                <Grid item>{renderItem("Xlog", JSON.stringify(nodePatroni.xlog))}</Grid>
            </Grid>
        </Grid>
    )

    function renderItem(name: string, value: string) {
        return (
            <div style={{ padding: '5px' }}>
                <b>{name}:</b>{` `}<span>{value}</span>
            </div>
        )
    }
}
