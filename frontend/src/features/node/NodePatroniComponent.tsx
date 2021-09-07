import {useAppDispatch, useAppSelector} from "../../app/hooks";
import {useEffect} from "react";
import {Grid, Paper} from "@material-ui/core";
import {incrementAsync, selectNodePatroniData} from "./NodePatroniSlice";
import css from "./NodePatroni.module.css"
import classNames from "classnames";

export function NodePatroniComponent() {
    const dispatch = useAppDispatch();
    const node = `p4-fr-ppl-1`;
    useEffect(() => { dispatch(incrementAsync(node)) }, [dispatch, node])

    const nodePatroni = useAppSelector(selectNodePatroniData)
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
