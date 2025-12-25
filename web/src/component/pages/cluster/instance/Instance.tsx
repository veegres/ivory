import {Box, Divider} from "@mui/material"

import {useRouterClusterOverview} from "../../../../api/cluster/hook"
import {SxPropsMap} from "../../../../app/type"
import {getQueryConnection} from "../../../../app/utils"
import {useStore, useStoreAction} from "../../../../provider/StoreProvider"
import {AlertCentered} from "../../../view/box/AlertCentered"
import {PageMainBox} from "../../../view/box/PageMainBox"
import {InstanceInfo} from "./InstanceInfo"
import {InstanceMain} from "./InstanceMain"

const SX: SxPropsMap = {
    content: {display: "flex", gap: 3},
}

export function Instance() {
    const {setInstanceBody} = useStoreAction
    const instance = useStore(s => s.instance)
    const activeCluster = useStore(s => s.activeCluster)
    const activeInstanceName = useStore(s => s.activeInstance[activeCluster?.cluster.name ?? ""])
    const activeClusterTab = useStore(s => s.activeClusterTab)
    const isClusterOverviewOpen = !!activeCluster && activeClusterTab === 0

    const overview = useRouterClusterOverview(activeCluster?.cluster.name, false)
    const activeInstance = overview.data && overview.data.instances[activeInstanceName ?? ""]
    return (
        <PageMainBox withPadding visible={isClusterOverviewOpen}>
            {renderContent()}
        </PageMainBox>
    )

    function renderContent() {
        if (!activeInstanceName || !activeCluster) return <AlertCentered text={"Please, select an instance to see the information!"}/>
        if (!activeInstance) return <AlertCentered text={"There is not enough information about the instance!"} severity={"warning"}/>

        const connection = getQueryConnection(activeCluster.cluster, activeInstance.database)

        return (
            <Box sx={SX.content}>
                <InstanceInfo instance={activeInstance} tab={instance.body} onTab={setInstanceBody} connection={connection}/>
                <Divider orientation={"vertical"} flexItem/>
                <InstanceMain tab={instance.body} database={activeInstance.database}/>
            </Box>
        )
    }
}
