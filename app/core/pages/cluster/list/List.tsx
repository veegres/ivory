import {Box} from "@mui/material"

import {useRouterClusterList} from "../../../../features/cluster/api/ClusterHook"
import {ClusterDeploy} from "../../../../features/cluster/component/ClusterDeploy"
import {ClusterDetect} from "../../../../features/cluster/component/ClusterDetect"
import {Feature} from "../../../../features/Feature"
import {ManageAccess} from "../../../../features/management/component/ManageAccess"
import {ErrorSmart} from "../../../../shared/component/box/ErrorSmart"
import {PageMainBox} from "../../../../shared/component/box/PageMainBox"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {KeeperPluginOptions, SxPropsFormatter} from "../../../../shared/helper/HelperUtils"
import {useStore} from "../../../../shared/provider/StoreProvider"
import {ListKeepers} from "./ListKeepers"
import {ListTable} from "./ListTable"
import {ListTableRefresher} from "./ListTableRefresher"
import {ListTags} from "./ListTags"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 0.5},
    filter: {display: "flex", flexDirection: "column", gap: "4px"},
    keepers: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 1, padding: "0 36px"},
    actions: {display: "flex", alignItems: "center", gap: "5px"},
}

export function List() {
    const clusters = useRouterClusterList()
    const keeper = useStore(s => s.activeClusterKeeperPlugin)

    return (
        <Box sx={SX.box}>
            <Box sx={[SX.filter, SxPropsFormatter.style.pageMargin]}>
                <Box sx={SX.keepers}>
                    <ListKeepers/>
                    <Box sx={SX.actions}>
                        <ClusterDeploy size={26} keeper={keeper} database={KeeperPluginOptions[keeper].dbPlugin}/>
                        <ClusterDetect size={26} keeper={keeper} database={KeeperPluginOptions[keeper].dbPlugin}/>
                        <ListTableRefresher size={26} width={55}/>
                    </Box>
                </Box>
                <ManageAccess feature={Feature.ViewTagList}>
                    <ListTags/>
                </ManageAccess>
            </Box>
            <PageMainBox>
                {clusters.error ? <ErrorSmart error={clusters.error}/> : (
                    <ListTable list={Object.values(clusters.data ?? [])} fetching={clusters.isFetching} pending={clusters.isPending}/>
                )}
            </PageMainBox>
        </Box>
    )
}
