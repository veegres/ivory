import {Box, Divider} from "@mui/material"

import {SxPropsMap} from "../../../../app/type"
import {getQueryConnection} from "../../../../app/utils"
import {useStore, useStoreAction} from "../../../../provider/StoreProvider"
import {InfoAlert} from "../../../view/box/InfoAlert"
import {PageMainBox} from "../../../view/box/PageMainBox"
import {InstanceInfo} from "./InstanceInfo"
import {InstanceMain} from "./InstanceMain"

const SX: SxPropsMap = {
    content: {display: "flex", gap: 3},
}

export function Instance() {
    const activeInstance = useStore(s => s.getActiveInstance())
    const activeCluster = useStore(s => s.activeCluster)
    const isClusterOverviewOpen = useStore(s => s.isClusterOverviewOpen)
    const instance = useStore(s => s.instance)
    const {setInstanceBody} = useStoreAction

    return (
        <PageMainBox withPadding visible={isClusterOverviewOpen()}>
            {renderContent()}
        </PageMainBox>
    )

    function renderContent() {
        if (!activeInstance || !activeCluster) return <InfoAlert text={"Please, select an instance to see the information!"}/>

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
