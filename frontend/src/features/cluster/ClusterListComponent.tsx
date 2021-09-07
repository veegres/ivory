import {useEffect} from "react";
import {useAppDispatch, useAppSelector} from "../../app/hooks";
import {incrementAsync, selectClusterListData} from "./ClusterListSlice";
import {Chip, Grid, Table, TableBody, TableCell, TableHead, TableRow} from "@material-ui/core";

export function ClusterListComponent() {
    const dispatch = useAppDispatch();
    useEffect(() => { dispatch(incrementAsync()) }, [dispatch])

    const clusterList = useAppSelector(selectClusterListData)
    if (clusterList.length === 0) return null;
    return (
        <Table>
            <TableHead>
                <TableRow>
                    <TableCell>Name</TableCell>
                    <TableCell>Nodes</TableCell>
                </TableRow>
            </TableHead>
            <TableBody>
                {clusterList.map(cluster => (
                    <TableRow key={cluster.name}>
                        <TableCell>{cluster.name}</TableCell>
                        <TableCell>
                            <Grid container spacing={1}>
                                {cluster.nodes.map(node => (
                                    <Grid item key={node}>
                                        <Chip variant="outlined" color="primary" clickable label={node} />
                                    </Grid>
                                ))}
                            </Grid>
                        </TableCell>
                    </TableRow>
                ))}
            </TableBody>
        </Table>
    )
}
