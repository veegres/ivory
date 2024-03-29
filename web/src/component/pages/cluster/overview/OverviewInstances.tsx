import {Table, TableBody, TableCell, TableHead, TableRow} from "@mui/material";
import {useIsFetching, useQueryClient} from "@tanstack/react-query";
import {TableCellLoader} from "../../../view/table/TableCellLoader";
import {TabProps} from "./Overview";
import {SxPropsMap} from "../../../../type/common";
import {OverviewInstancesRow} from "./OverviewInstancesRow";
import {RefreshIconButton} from "../../../view/button/IconButtons";

const SX: SxPropsMap = {
    table: {"tr:last-child td": {border: 0}, "tr th, td": {padding: "2px 5px"}, tableLayout: "fixed"},
    buttonCell: {width: "160px"},
    warning: {display: "flex", justifyContent: "center"},
}

export function OverviewInstances({info}: TabProps) {
    const {cluster, combinedInstanceMap} = info
    const key = {queryKey: ["instance/overview", cluster.name]}
    const queryClient = useQueryClient();
    const instanceMapFetching = useIsFetching(key)
    const candidates = Object.values(combinedInstanceMap)
        .filter(sidecar => sidecar.role === "replica")
        .map(instance => instance.sidecar)

    return (
        <Table size={"small"} sx={SX.table}>
            <TableHead>
                <TableRow>
                    <TableCell width={"44px"}/>
                    <TableCell width={"40px"}/>
                    <TableCell width={"110px"}>Role</TableCell>
                    <TableCell width={"15%"}>Sidecar</TableCell>
                    <TableCell width={"15%"}>Postgres</TableCell>
                    <TableCell width={"100px"}>State</TableCell>
                    <TableCell/>
                    <TableCellLoader sx={SX.buttonCell} isFetching={instanceMapFetching > 0}>
                        <RefreshIconButton onClick={() => queryClient.refetchQueries(key)}/>
                    </TableCellLoader>
                </TableRow>
            </TableHead>
            <TableBody>
                {Object.entries(combinedInstanceMap).map(([key, element]) => (
                    <OverviewInstancesRow key={key} domain={key} instance={element} cluster={cluster} candidates={candidates}/>
                ))}
            </TableBody>
        </Table>
    )
}
