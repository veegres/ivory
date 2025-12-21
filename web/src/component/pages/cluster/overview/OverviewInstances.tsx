import {Box, Table, TableBody, TableCell, TableHead, TableRow} from "@mui/material"

import {useRouterClusterOverview} from "../../../../api/cluster/hook"
import {Cluster, Instance} from "../../../../api/cluster/type"
import {SxPropsMap} from "../../../../app/type"
import {getActiveInstance, getDomain, getIsActiveInstance} from "../../../../app/utils"
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
    instances?: { [domain: string]: Instance },
}

export function OverviewInstances(props: Props) {
    const {cluster} = props
    const instances = props.instances ?? cluster.unknownInstances
    const activeInstances = useStore(s => s.activeInstance)
    const activeInstance = getActiveInstance(activeInstances, cluster)
    const overview = useRouterClusterOverview(cluster.name)
    const candidates = Object.values(instances)
        .filter(sidecar => sidecar.role === "replica")
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
                        <OverviewInstancesRow key={key} checked={getIsActiveInstance(key, activeInstance)} instance={element} cluster={cluster} candidates={candidates}/>
                    ))}
                </TableBody>
            </Table>
        </Box>
    )

    function renderCheckedInstance() {
        if (!activeInstance) return
        const domain = getDomain(activeInstance.sidecar)
        if (instances[domain]) return
        const checked = getIsActiveInstance(domain, activeInstance)
        return (
            <OverviewInstancesRow checked={checked} instance={activeInstance} cluster={cluster} candidates={candidates} error={true}/>
        )
    }
}
