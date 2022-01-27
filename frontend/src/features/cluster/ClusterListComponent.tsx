import {Chip, Grid, Table, TableBody, TableCell, TableHead, TableRow} from "@material-ui/core";
import { useQuery } from "react-query";
import {getClusterList} from "../../app/api";

export function ClusterListComponent() {
    const { data: clusterList } = useQuery('cluster/list', getClusterList)
    console.log(clusterList);
    if (!clusterList || clusterList.length === 0) return null;
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
