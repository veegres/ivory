import {Box} from "@mui/material"
import {useMemo, useState} from "react"

import {Cluster} from "../../../../features/cluster/api/ClusterType"
import {ClusterDeploy} from "../../../../features/cluster/component/ClusterDeploy"
import {ClusterDetect} from "../../../../features/cluster/component/ClusterDetect"
import {AlertCentered} from "../../../../shared/component/box/AlertCentered"
import {ActionsLoader} from "../../../../shared/component/progress/ActionsLoader"
import {SkeletonRows} from "../../../../shared/component/progress/SkeletonRows"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {KeeperPluginOptions, SxPropsFormatter} from "../../../../shared/helper/HelperUtils"
import {useStore} from "../../../../shared/provider/StoreProvider"
import scroll from "../../../../shared/style/scroll.module.css"
import {ListClusterAdd} from "./ListClusterAdd"
import {ListEmptyInfo} from "./ListEmptyInfo"
import {ListRow} from "./ListRow"
import {ListRowNew} from "./ListRowNew"
import {ListTableRefresher} from "./ListTableRefresher"

const SX: SxPropsMap = {
    box: {overflowY: "scroll"},
    head: {
        position: "sticky", top: 0, zIndex: 2, display: "flex", alignItems: "center", gap: 1,
        padding: "5px", fontFamily: "monospace", letterSpacing: 1,
        bgcolor: "background.paper", borderBottom: 1, borderColor: "divider",
    },
    headName: {display: {xs: "none", md: "block"}, flex: {md: "0 0 206px"}, padding: "0px 5px"},
    headClusters: {display: {xs: "block", md: "none"}, padding: "0px 5px"},
    headNodes: {display: {xs: "none", md: "block"}, padding: "0px 5px"},
    empty: {padding: "5px"},
}

type Props = {
    list: Cluster[],
    pending: boolean,
    fetching: boolean,
}

export function ListTable(props: Props) {
    const keeper = useStore(s => s.activeClusterKeeperPlugin)
    const activeCluster = useStore(s => s.activeCluster)
    const search = useStore(s => s.searchCluster)
    const {list, fetching, pending} = props
    const [showNewElement, setShowNewElement] = useState(false)
    const [editNode, setEditNode] = useState("")

    const rows = useMemo(() => list.filter((c) => c.name.includes(search)), [list, search])

    return (
        <Box sx={[SX.box, {maxHeight: activeCluster ? "25vh" : "60vh"}]} className={scroll.tiny}>
            <Box sx={[SX.head, SxPropsFormatter.style.paper]}>
                <Box sx={SX.headName}>Name</Box>
                <Box sx={SX.headClusters}>Clusters</Box>
                <ActionsLoader label={<Box sx={SX.headNodes}>Nodes</Box>} loading={fetching && !pending}>
                    <ListTableRefresher/>
                    <ClusterDeploy keeper={keeper} database={KeeperPluginOptions[keeper].dbPlugin}/>
                    <ClusterDetect keeper={keeper} database={KeeperPluginOptions[keeper].dbPlugin}/>
                    <ListClusterAdd onClick={() => setShowNewElement(true)} disabled={showNewElement}/>
                </ActionsLoader>
            </Box>
            <SkeletonRows isLoading={pending} height={32}>
                <ListRowNew show={showNewElement} close={() => setShowNewElement(false)}/>
                {renderRemovedRow()}
                {renderRows()}
                {renderEmpty()}
            </SkeletonRows>
        </Box>
    )

    function renderRemovedRow() {
        if (!activeCluster) return
        if (rows.some(e => e.name === activeCluster.name)) return
        return (
            <ListRow cluster={activeCluster} editable={false}/>
        )
    }

    function renderRows() {
        return rows.map((cluster) => {
            const editable = cluster.name === editNode
            const toggle = () => setEditNode(editable ? "" : cluster.name)
            return (
                <ListRow key={cluster.name} cluster={cluster} editable={editable} toggle={toggle}/>
            )
        })
    }

    function renderEmpty() {
        if (pending || showNewElement || rows.length || activeCluster) return
        return (
            <Box sx={SX.empty}>
                {search ? (
                    <AlertCentered text={"There are no clusters that match your filter"}/>
                ) : (
                    <ListEmptyInfo
                        onAddManually={() => setShowNewElement(true)}
                        disabledAddManually={showNewElement}
                    />
                )}
            </Box>
        )
    }
}
