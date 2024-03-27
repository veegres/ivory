import {Box, Divider} from "@mui/material";
import {useStore, useStoreAction} from "../../../../provider/StoreProvider";
import {InfoAlert} from "../../../view/box/InfoAlert";
import {PageMainBox} from "../../../view/box/PageMainBox";
import {SxPropsMap} from "../../../../type/general";
import {InstanceMain} from "./InstanceMain";
import {InstanceInfo} from "./InstanceInfo";

const SX: SxPropsMap = {
    content: {display: "flex", gap: 3},
}

export function Instance() {
    const {activeInstance, isClusterOverviewOpen, instance} = useStore()
    const {setInstanceBody} = useStoreAction()

    return (
        <PageMainBox withPadding visible={isClusterOverviewOpen()}>
            {renderContent()}
        </PageMainBox>
    )

    function renderContent() {
        if (!activeInstance) return <InfoAlert text={"Please, select an instance to see the information!"}/>

        return (
            <Box sx={SX.content}>
                <InstanceInfo instance={activeInstance} tab={instance.body} onTab={setInstanceBody}/>
                <Divider orientation={"vertical"} flexItem/>
                <InstanceMain tab={instance.body} database={activeInstance.database}/>
            </Box>
        )
    }
}
