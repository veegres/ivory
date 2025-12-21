import {Box, Divider} from "@mui/material"

import {SxPropsMap} from "../../../../app/type"
import {getActiveInstance, getQueryConnection} from "../../../../app/utils"
import {useStore, useStoreAction} from "../../../../provider/StoreProvider"
import {AlertCentered} from "../../../view/box/AlertCentered"
import {PageMainBox} from "../../../view/box/PageMainBox"
import {InstanceInfo} from "./InstanceInfo"
import {InstanceMain} from "./InstanceMain"

const SX: SxPropsMap = {
    content: {display: "flex", gap: 3},
}

export function Instance() {
    const instance = useStore(s => s.instance)
    const activeInstances = useStore(s => s.activeInstance)
    const activeCluster = useStore(s => s.activeCluster)
    const activeClusterTab = useStore(s => s.activeClusterTab)
    const isClusterOverviewOpen = !!activeCluster && activeClusterTab === 0
    const activeInstance = getActiveInstance(activeInstances, activeCluster?.cluster)

    const {setInstanceBody} = useStoreAction

    return (
        <PageMainBox withPadding visible={isClusterOverviewOpen}>
            {renderContent()}
        </PageMainBox>
    )

    function renderContent() {
        if (!activeInstance || !activeCluster) return <AlertCentered text={"Please, select an instance to see the information!"}/>

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
