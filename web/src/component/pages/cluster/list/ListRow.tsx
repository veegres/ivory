import {ErrorOutlineRounded, WarningAmberRounded} from "@mui/icons-material"
import {Box, TableRow, Tooltip} from "@mui/material"
import {useEffect, useMemo, useRef, useState} from "react"

import {useRouterClusterOverview} from "../../../../api/cluster/hook"
import {Cluster} from "../../../../api/cluster/type"
import {ColorsMap, SxPropsMap} from "../../../../app/type"
import {getDetectionItems, getDomains, getNodeConnections, NodeColor, SxPropsFormatter} from "../../../../app/utils"
import {useStore, useStoreAction} from "../../../../provider/StoreProvider"
import {InfoColorBoxList} from "../../../view/box/InfoColorBoxList"
import {AutoRefreshIconButton, RefreshIconButton} from "../../../view/button/IconButtons"
import {DynamicInputs} from "../../../view/input/DynamicInputs"
import {ListCell} from "./ListCell"
import {ListCellChip} from "./ListCellChip"
import {ListCellRead} from "./ListCellRead"
import {ListCellUpdate} from "./ListCellUpdate"

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

    const activeCluster = useStore(s => s.activeCluster)
    const {setWarnings, setCluster} = useStoreAction
    const active = !!activeCluster && cluster.name === activeCluster.cluster.name

    const ref = useRef<HTMLTableRowElement | null>(null)
    const [stateNodes, setStateNodes] = useState(getDomains(cluster.nodes))

    const overview = useRouterClusterOverview(cluster.name)
    const nodes = overview.data?.nodes ?? cluster.nodesOverview
    const detectedDomain = overview.data?.detectedDomain
    const [warning, colors] = useMemo(handleMemoNodes, [nodes])
    const detectionItems = useMemo(() => getDetectionItems(nodes, detectedDomain, activeCluster?.manualKeeper), [nodes, detectedDomain, activeCluster?.manualKeeper])

    useEffect(handleEffectWarningsUpdate, [cluster.name, setWarnings, warning])
    useEffect(handleEffectScroll, [active])
    useEffect(handleEffectNodesUpdate, [cluster.nodes])
    useEffect(handleEffectActiveClusterUpdate, [active, cluster, setCluster, warning])

    return (
        <TableRow sx={[active && SxPropsFormatter.style.bgImageSelected, !toggle && SxPropsFormatter.style.bgImageError]} ref={ref}>
            <ListCell width={"220px"}>
                <ListCellChip name={cluster.name} active={active} onClick={handleClick} renderRefresh={renderRefresh()}/>
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
                cluster={{...cluster, nodes: getNodeConnections(cluster.keeperType, cluster.dbType, stateNodes)}}
                toggle={toggle}
                onUpdate={overview.refetch}
                onClose={() => setStateNodes(getDomains(cluster.nodes))}
            />
        )
    }

    function renderRefresh() {
        return !active || !activeCluster?.manualKeeper ? (
            <AutoRefreshIconButton
                tooltip={renderChipTooltip()}
                placement={"right-start"}
                arrow={true}
                loading={overview.isFetching}
                onClick={() => overview.refetch()}
            />
        ) : (
            <RefreshIconButton
                tooltip={renderChipTooltip()}
                placement={"right-start"}
                arrow={true}
                loading={overview.isFetching}
                onClick={() => overview.refetch()}
            />
        )
    }

    function renderChipTooltip() {
        return <InfoColorBoxList items={detectionItems} label={"Cluster Detection"}/>
    }

    function handleMemoNodes(): [boolean, ColorsMap] {
        const getColors = () => Object.entries(nodes).reduce(
            (map, [domain, node]) => {
                map[domain] = NodeColor[node.keeper.role].label
                return map
            },
            {} as ColorsMap
        )
        const getWarning = () => {
            for (const key in nodes) {
                const node = nodes[key]
                if (node.warnings.length > 0) return true
            }
            return false
        }
        return [getWarning(), getColors()]
    }


    function handleEffectWarningsUpdate() {
        setWarnings(cluster.name, warning)
        return () => {
            if (warning) setWarnings(cluster.name, false)
        }
    }


    // NOTE: we could potentially take the data from tanstack.query, but then we couldn't
    //  see the selected cluster which is not on the list
    function handleEffectActiveClusterUpdate() {
        if (!active) return
        setCluster({cluster, warning})
    }

    function handleEffectScroll() {
        if (active) ref.current?.scrollIntoView({block: "nearest"})
    }

    function handleEffectNodesUpdate() {
        setStateNodes(getDomains(cluster.nodes))
    }

    function handleClick() {
        if (active) {
            setCluster(undefined)
        } else {
            setCluster({cluster, warning})
        }
    }
}
