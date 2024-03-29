import {Table, TableBody, TableCell, TableHead, TableRow} from "@mui/material";
import {useIsFetching, useQueryClient} from "@tanstack/react-query";
import {TableCellLoader} from "../../../view/table/TableCellLoader";
import {SxPropsMap} from "../../../../type/general";
import {OverviewInstancesRow} from "./OverviewInstancesRow";
import {RefreshIconButton} from "../../../view/button/IconButtons";
import {ActiveCluster} from "../../../../type/cluster";
import {InstanceWeb} from "../../../../type/instance";
import {getDomain} from "../../../../app/utils";
import {InstanceApi} from "../../../../app/api";

const SX: SxPropsMap = {
    table: {"tr:last-child td": {border: 0}, "tr th, td": {padding: "2px 5px"}, tableLayout: "fixed"},
    buttonCell: {width: "160px"},
    warning: {display: "flex", justifyContent: "center"},
}

type Props = {
    info: ActiveCluster,
    activeInstance?: InstanceWeb,
}

export function OverviewInstances(props: Props) {
    const {activeInstance, info} = props
    const {cluster, combinedInstanceMap} = info
    const key = {queryKey: InstanceApi.overview.key(cluster.name)}
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
                {renderCheckedInstance()}
                {Object.entries(combinedInstanceMap).map(([key, element]) => (
                    <OverviewInstancesRow key={key} domain={key} instance={element} cluster={cluster} candidates={candidates}/>
                ))}
            </TableBody>
        </Table>
    )

    function renderCheckedInstance() {
        if (!activeInstance) return
        const domain = getDomain(activeInstance.sidecar)
        if (combinedInstanceMap[domain]) return
        return (
            <OverviewInstancesRow domain={domain} instance={activeInstance} cluster={cluster} candidates={candidates} error={true}/>
        )
    }
}
