import {ErrorOutlineRounded, WarningAmberRounded} from "@mui/icons-material"
import {Box, TableRow, Tooltip} from "@mui/material"
import {useEffect, useMemo, useRef, useState} from "react"

import {useRouterClusterOverview} from "../../../../api/cluster/hook"
import {Cluster} from "../../../../api/cluster/type"
import {ColorsMap, SxPropsMap} from "../../../../app/type"
import {getDomain, getDomains, getNodeConnections, NodeColor, SxPropsFormatter} from "../../../../app/utils"
import {useStore, useStoreAction} from "../../../../provider/StoreProvider"
import {DynamicInputs} from "../../../view/input/DynamicInputs"
import {ListCell} from "./ListCell"
import {ListCellRead} from "./ListCellRead"
import {ListCellUpdate} from "./ListCellUpdate"
import {ListRowName} from "./ListRowName"

const SX: SxPropsMap = {
    actions: {display: "flex", justifyContent: "flex-end", alignItems: "center", gap: 1, minHeight: "32px"},
    warning: {padding: "3px 0"},
    removed: {color: "text.secondary", fontSize: "12px", textAlign: "center", marginTop: "5px"},
}

type Props = {
    cluster: Cluster,
    editable: boolean,
    toggle?: () => void,
}

export function ListRow(props: Props) {
    const {cluster, editable, toggle} = props
    const ref = useRef<HTMLTableRowElement | null>(null)

    const {setWarnings} = useStoreAction
    const activeCluster = useStore(s => s.activeCluster)
    const active = !!activeCluster && cluster.name === activeCluster.name

    const [stateNodes, setStateNodes] = useState(getDomains(cluster.nodes, !editable))

    const overview = useRouterClusterOverview(cluster.name)
    const [warning, colors] = useMemo(handleMemoNodes, [overview.data?.nodes, editable])

    useEffect(handleEffectNodesUpdate, [cluster.nodes, editable])
    useEffect(handleEffectWarningsUpdate, [cluster.name, warning, setWarnings])
    useEffect(handleEffectScroll, [active])

    return (
        <TableRow sx={[active && SxPropsFormatter.style.bgImageSelected, !toggle && SxPropsFormatter.style.bgImageError]} ref={ref}>
            <ListCell width={"220px"}>
                <ListRowName cluster={cluster} active={active} loading={overview.isFetching} refresh={overview.refetch}/>
            </ListCell>
            <ListCell>
                <DynamicInputs
                    inputs={stateNodes}
                    colors={colors}
                    editable={editable}
                    placeholder={"Node "}
                    onChange={n => setStateNodes(n)}
                />
            </ListCell>
            <ListCell width={"130px"}>
                <Box sx={SX.actions}>
                    {warning && !overview.error && !overview.isFetching && (
                        <Tooltip title={"Issues detected — select a cluster to view details"} placement={"top"}>
                            <WarningAmberRounded color={"warning"}/>
                        </Tooltip>
                    )}
                    {overview.error && (
                        <Tooltip title={overview.error.message} placement={"top"}>
                            <ErrorOutlineRounded color={"error"}/>
                        </Tooltip>
                    )}
                    {!toggle && (
                        <Tooltip title={"This cluster isn't in the list — it appears here because it was manually selected. Uncheck it to remove it."} placement={"top"}>
                            <ErrorOutlineRounded color={"secondary"}/>
                        </Tooltip>
                    )}
                    {renderButtons()}
                </Box>
            </ListCell>
        </TableRow>
    )

    function renderButtons() {
        if (!toggle) return
        return !editable ? (
            <ListCellRead name={cluster.name} toggle={toggle}/>
        ) : (
            <ListCellUpdate
                cluster={{...cluster, nodes: getNodeConnections(stateNodes)}}
                toggle={toggle}
                onUpdate={overview.refetch}
                onClose={() => setStateNodes(getDomains(cluster.nodes))}
            />
        )
    }

    function handleMemoNodes(): [boolean, ColorsMap] {
        let warning = false
        const colors = Object.values(overview.data?.nodes ?? {}).reduce(
            (map, node) => {
                if (node.warnings.length > 0) warning = true
                const d = getDomain(node.config, !editable)
                map[d] = NodeColor[node.keeper.role].label
                return map
            },
            {} as ColorsMap
        )
        return [warning, colors]
    }

    function handleEffectWarningsUpdate() {
        setWarnings(cluster.name, warning)
        return () => {
            if (warning) setWarnings(cluster.name, false)
        }
    }

    function handleEffectScroll() {
        if (active) ref.current?.scrollIntoView({block: "nearest"})
    }

    function handleEffectNodesUpdate() {
        setStateNodes(getDomains(cluster.nodes, !editable))
    }
}
