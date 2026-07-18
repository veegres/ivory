import {ErrorOutlineRounded, WarningAmberRounded} from "@mui/icons-material"
import {Box, Tooltip} from "@mui/material"
import {useEffect, useMemo, useRef, useState} from "react"

import {useRouterClusterOverview} from "../../../../features/cluster/api/ClusterHook"
import {Cluster} from "../../../../features/cluster/api/ClusterType"
import {ColorsMap, SxPropsMap} from "../../../../shared/helper/HelperType"
import {getDomain, getDomains, getNodeConfigs, NodeColor, SxPropsFormatter} from "../../../../shared/helper/HelperUtils"
import {useStore, useStoreAction} from "../../../../shared/provider/StoreProvider"
import {ListCellRead} from "./ListCellRead"
import {ListCellUpdate} from "./ListCellUpdate"
import {ListNodeInput} from "./ListNodeInput"
import {ListRowLayout} from "./ListRowLayout"
import {ListRowName} from "./ListRowName"

// NOTE: status mirrors the 32px footprint of IconButton so the warning/error
// icons line up with the edit/delete buttons next to them; the glyph is a bit
// bigger than the buttons' 18px because these outline icons read smaller at
// the same font size
const SX: SxPropsMap = {
    actions: {display: "flex", justifyContent: "flex-end", alignItems: "center", minHeight: "32px"},
    status: {width: "32px", height: "32px", display: "flex", alignItems: "center", justifyContent: "center", cursor: "pointer"},
    statusIcon: {fontSize: 23},
}

type Props = {
    cluster: Cluster,
    editable: boolean,
    toggle?: () => void,
}

export function ListRow(props: Props) {
    const {cluster, editable, toggle} = props
    const ref = useRef<HTMLDivElement | null>(null)

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
        <ListRowLayout
            sx={[active && SxPropsFormatter.style.bgImageSelected, !toggle && SxPropsFormatter.style.bgImageError]}
            ref={ref}
            renderName={renderName()}
            renderNodes={renderNodes()}
            renderActions={renderActions()}
        />
    )

    function renderName() {
        return (
            <ListRowName cluster={cluster} active={active} loading={overview.isFetching} refresh={overview.refetch}/>
        )
    }

    function renderNodes() {
        return (
            <ListNodeInput
                inputs={stateNodes}
                colors={colors}
                editable={editable}
                onChange={n => setStateNodes(n)}
            />
        )
    }

    function renderActions() {
        return (
            <Box sx={SX.actions}>
                {warning && !overview.error && !overview.isFetching && (
                    <Tooltip title={"Issues detected — select a cluster to view details"} placement={"top"}>
                        <Box sx={SX.status}><WarningAmberRounded sx={SX.statusIcon} color={"warning"}/></Box>
                    </Tooltip>
                )}
                {overview.error && (
                    <Tooltip title={overview.error.message} placement={"top"}>
                        <Box sx={SX.status}><ErrorOutlineRounded sx={SX.statusIcon} color={"error"}/></Box>
                    </Tooltip>
                )}
                {!toggle && (
                    <Tooltip title={"This cluster isn't in the list — it appears here because it was manually selected. Uncheck it to remove it."} placement={"top"}>
                        <Box sx={SX.status}><ErrorOutlineRounded sx={SX.statusIcon} color={"secondary"}/></Box>
                    </Tooltip>
                )}
                {renderButtons()}
            </Box>
        )
    }

    function renderButtons() {
        if (!toggle) return
        return !editable ? (
            <ListCellRead name={cluster.name} toggle={toggle}/>
        ) : (
            <ListCellUpdate
                cluster={{...cluster, nodes: getNodeConfigs(stateNodes)}}
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
