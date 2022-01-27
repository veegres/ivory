import {Button, Grid, Table, TableBody, TableCell, TableHead, TableRow} from "@material-ui/core";
import {OpenInNew} from "@material-ui/icons";
import { useQuery } from "react-query";
import {getNodeCluster} from "../../app/api";

export function ClusterInfoComponent() {
    const node = `P4-IO-CHAT-10`;
    const { data: members } = useQuery(['node/cluster', node], () => getNodeCluster(node))
    if (!members || members.length === 0) return null;

    return (
        <Table>
            <TableHead>
                <TableRow>
                    <TableCell>Node</TableCell>
                    <TableCell>Role</TableCell>
                    <TableCell>Lag</TableCell>
                    <TableCell />
                </TableRow>
            </TableHead>
            <TableBody>
                {members.map((node) => (
                    <TableRow key={node.host}>
                        <TableCell>{node.name}</TableCell>
                        <TableCell>{node.role.toUpperCase()}</TableCell>
                        <TableCell>{node.lag}</TableCell>
                        <TableCell align="right">
                            <Grid container spacing={3} justifyContent="flex-end" alignItems="center">
                                <Grid item>{node.role === "leader" ? <Button color="secondary">Switchover</Button> : null}</Grid>
                                <Grid item><Button color="primary">Reinit</Button></Grid>
                                <Grid item><Button href={node.api_url} target="_blank"><OpenInNew /></Button></Grid>
                            </Grid>
                        </TableCell>
                    </TableRow>
                ))}
            </TableBody>
        </Table>
    )
}
