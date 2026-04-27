import {Box, Table, TableBody, TableCell, TableHead, TableRow} from "@mui/material"

import {useRouterClusterOverview} from "../../../../api/cluster/hook"
import {Cluster, NodeOverview} from "../../../../api/cluster/type"
import {SxPropsMap} from "../../../../app/type"
import {useStore} from "../../../../provider/StoreProvider"
import {RefreshIconButton} from "../../../view/button/IconButtons"
import {TableCellLoader} from "../../../view/table/TableCellLoader"
import {OverviewNodesFixAuto} from "./OverviewNodesFixAuto"
import {OverviewNodesRow} from "./OverviewNodesRow"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 2},
    table: {"tr:last-child td": {border: 0}, "tr th, td": {padding: "2px 5px"}, tableLayout: "fixed"},
    warning: {display: "flex", justifyContent: "center"},
}

type Props = {
    cluster: Cluster,
    nodes?: NodeOverview,
}

export function OverviewNodes(props: Props) {
    const {cluster} = props
    const nodes = props.nodes ?? cluster.nodesOverview
    const activeNode = useStore(s => s.activeNode[cluster.name])
    const overview = useRouterClusterOverview(cluster.name, false)
    const candidates = Object.values(nodes)
        .filter(node => !!node)
        .filter(node => node.keeper.role === "replica")
        .map(node => node.connection)

    return (
        <Box sx={SX.box}>
            <Table size={"small"} sx={SX.table}>
                <TableHead>
                    <TableRow>
                        <TableCell sx={{width: "44px"}}/>
                        <TableCell sx={{width: "40px"}}/>
                        <TableCell sx={{width: "110px"}}>Role</TableCell>
                        <TableCell sx={{width: "15%"}}>Host</TableCell>
                        <TableCell sx={{width: "60px"}}>SSH</TableCell>
                        <TableCell sx={{width: "70px"}}>Keeper</TableCell>
                        <TableCell sx={{width: "80px"}}>Database</TableCell>
                        <TableCell sx={{width: "10%"}}>State</TableCell>
                        <TableCellLoader sx={{width: "75%"}} loading={false}>
                            <OverviewNodesFixAuto name={cluster.name}/>
                            <RefreshIconButton
                                onClick={() => overview.refetch()}
                                loading={overview.isFetching}
                            />
                        </TableCellLoader>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {renderCheckedNode()}
                    {Object.entries(nodes).map(([key, element]) => (
                        <OverviewNodesRow
                            key={key}
                            name={key}
                            checked={key === activeNode}
                            node={element}
                            cluster={cluster}
                            candidates={candidates}
                        />
                    ))}
                </TableBody>
            </Table>
        </Box>
    )

    function renderCheckedNode() {
        if (!activeNode) return
        if (nodes.hasOwnProperty(activeNode)) return
        return (
            <OverviewNodesRow name={activeNode} checked={true} cluster={cluster} candidates={candidates} error={true}/>
        )
    }
}
