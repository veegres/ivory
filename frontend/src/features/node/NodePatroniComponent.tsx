import {Grid, Paper} from "@material-ui/core";
import css from "./NodePatroni.module.css"
import classNames from "classnames";
import {getNodePatroni} from "../../app/api";
import { useQuery } from "react-query";

export function NodePatroniComponent() {
    const node = `P4-IO-CHAT-10`;
    const { data: nodePatroni } = useQuery(['node/patroni', node], () => getNodePatroni(node))
    if (!nodePatroni) return null;

    const roleClassName = classNames(css.role, { [css.green]: nodePatroni.role === "master" })

    return (
        <Paper className={css.paper} elevation={3}>
            <Grid container direction="row" spacing={3}>
                <Grid className={roleClassName} item xs={2} container alignContent="center" justifyContent="center">
                    <Grid item>{nodePatroni.role.toUpperCase()}</Grid>
                </Grid>
                <Grid item xs={10} container direction="column" spacing={1}>
                    <Grid item>{renderItem("Node", node)}</Grid>
                    <Grid item>{renderItem("State", nodePatroni.state)}</Grid>
                    <Grid item>{renderItem("Scope", nodePatroni.patroni.scope)}</Grid>
                    <Grid item>{renderItem("Timeline", nodePatroni.timeline.toString())}</Grid>
                    <Grid item>{renderItem("Xlog", JSON.stringify(nodePatroni.xlog))}</Grid>
                </Grid>
            </Grid>
        </Paper>
    )

    function renderItem(name: string, value: string) {
        return (
            <>
                <b>{name}:</b>{` `}<span>{value}</span>
            </>
        )
    }
}
