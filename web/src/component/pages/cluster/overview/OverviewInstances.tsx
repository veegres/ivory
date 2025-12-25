import {Box, Table, TableBody, TableCell, TableHead, TableRow} from "@mui/material"

import {useRouterClusterOverview} from "../../../../api/cluster/hook"
import {Cluster, InstanceOverview} from "../../../../api/cluster/type"
import {SxPropsMap} from "../../../../app/type"
import {useStore} from "../../../../provider/StoreProvider"
import {RefreshIconButton} from "../../../view/button/IconButtons"
import {TableCellLoader} from "../../../view/table/TableCellLoader"
import {OverviewInstancesFixAuto} from "./OverviewInstancesFixAuto"
import {OverviewInstancesRow} from "./OverviewInstancesRow"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 2},
    table: {"tr:last-child td": {border: 0}, "tr th, td": {padding: "2px 5px"}, tableLayout: "fixed"},
    buttonCell: {width: "160px"},
    warning: {display: "flex", justifyContent: "center"},
}

type Props = {
    cluster: Cluster,
    instances?: InstanceOverview,
}

export function OverviewInstances(props: Props) {
    const {cluster} = props
    const instances = props.instances ?? cluster.sidecarsOverview
    const activeInstance = useStore(s => s.activeInstance[cluster.name])
    const overview = useRouterClusterOverview(cluster.name, false)
    const candidates = Object.values(instances)
        .filter(instance => !!instance)
        .filter(instance => instance.role === "replica")
        .map(instance => instance.sidecar)

    return (
        <Box sx={SX.box}>
            <Table size={"small"} sx={SX.table}>
                <TableHead>
                    <TableRow>
                        <TableCell width={"44px"}/>
                        <TableCell width={"40px"}/>
                        <TableCell width={"110px"}>Role</TableCell>
                        <TableCell width={"15%"}>Sidecar</TableCell>
                        <TableCell width={"15%"}>Postgres</TableCell>
                        <TableCell width={"150px"}>State</TableCell>
                        <TableCell/>
                        <TableCellLoader sx={SX.buttonCell} isFetching={false}>
                            <OverviewInstancesFixAuto name={cluster.name}/>
                            <RefreshIconButton
                                onClick={() => overview.refetch()}
                                loading={overview.isFetching}
                            />
                        </TableCellLoader>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {renderCheckedInstance()}
                    {Object.entries(instances).map(([key, element]) => (
                        <OverviewInstancesRow
                            key={key}
                            name={key}
                            checked={key === activeInstance}
                            instance={element}
                            cluster={cluster}
                            candidates={candidates}
                        />
                    ))}
                </TableBody>
            </Table>
        </Box>
    )

    function renderCheckedInstance() {
        if (!activeInstance) return
        if (instances.hasOwnProperty(activeInstance)) return
        return (
            <OverviewInstancesRow name={activeInstance} checked={true} cluster={cluster} candidates={candidates} error={true}/>
        )
    }
}
