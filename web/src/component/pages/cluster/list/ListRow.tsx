import {ErrorOutlineRounded, WarningAmberRounded} from "@mui/icons-material"
import {Box, TableRow, Tooltip} from "@mui/material"
import {useEffect, useMemo, useRef, useState} from "react"

import {useRouterClusterOverview} from "../../../../api/cluster/hook"
import {Cluster} from "../../../../api/cluster/type"
import {ColorsMap, SxPropsMap} from "../../../../app/type"
import {
    getDetectionItems,
    getDomains,
    getKeepers, NodeColor,
    SxPropsFormatter,
} from "../../../../app/utils"
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
    const [stateNodes, setStateNodes] = useState(getDomains(cluster.keepers))

    const overview = useRouterClusterOverview(cluster.name)
    const nodes = overview.data?.nodes ?? cluster.keepersOverview
    const mainNode = overview.data?.mainNode
    const warning = useMemo(handleMemoWarning, [nodes])
    const colors = useMemo(handleMemoColors, [nodes])

    useEffect(handleEffectWarningsUpdate, [cluster.name, setWarnings, warning])
    useEffect(handleEffectScroll, [active])
    useEffect(handleEffectNodesUpdate, [cluster.keepers])
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
                        <Tooltip title={"Problems were detected, select the cluster to see it"} placement={"top"}>
                            <WarningAmberRounded color={"warning"}/>
                        </Tooltip>
                    )}
                    {overview.error && (
                        <Tooltip title={overview.error.message} placement={"top"}>
                            <ErrorOutlineRounded color={"error"}/>
                        </Tooltip>
                    )}
                    {!toggle && (
                        <Tooltip title={"This cluster is not in this list. It was selected that is why you see it. Just uncheck it."} placement={"top"}>
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
                cluster={{...cluster, keepers: getKeepers(stateNodes)}}
                toggle={toggle}
                onUpdate={overview.refetch}
                onClose={() => setStateNodes(getDomains(cluster.keepers))}
            />
        )
    }

    function renderRefresh() {
        return !active || !activeCluster?.detectBy ? (
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
        const detectBy = activeCluster?.detectBy
        const items = getDetectionItems(mainNode, detectBy)
        return <InfoColorBoxList items={items} label={"Cluster Detection"}/>
    }

    function handleMemoColors() {
        return Object.entries(nodes).reduce(
            (map, [domain, node]) => {
                map[domain] = NodeColor[node?.role ?? "unknown"].label
                return map
            },
            {} as ColorsMap
        )
    }

    function handleMemoWarning() {
        for (const key in nodes) {
            const node = nodes[key]
            if (!node) return true
        }
        return false
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
        setStateNodes(getDomains(cluster.keepers))
    }

    function handleClick() {
        if (active) {
            setCluster(undefined)
        } else {
            setCluster({cluster, warning})
        }
    }
}
