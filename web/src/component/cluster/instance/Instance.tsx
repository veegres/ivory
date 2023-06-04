import {Box, Divider} from "@mui/material";
import {useStore} from "../../../provider/StoreProvider";
import {InfoAlert} from "../../view/box/InfoAlert";
import {PageBox} from "../../view/box/PageBox";
import {SxPropsMap} from "../../../type/common";
import {useState} from "react";
import {InstanceMain} from "./InstanceMain";
import {InstanceInfo} from "./InstanceInfo";
import {InstanceTabType} from "../../../type/instance";

const SX: SxPropsMap = {
    content: {display: "flex", gap: 3},
}

export function Instance() {
    const {store: {activeInstance}, isClusterOverviewOpen} = useStore()
    const [tab, setTab] = useState(InstanceTabType.CHART)

    return (
        <PageBox withPadding visible={isClusterOverviewOpen()}>
            {renderContent()}
        </PageBox>
    )

    function renderContent() {
        if (!activeInstance) return <InfoAlert text={"Please, select a instance to see the information!"}/>

        return (
            <Box sx={SX.content}>
                <InstanceInfo instance={activeInstance} tab={tab} onTab={setTab}/>
                <Divider orientation={"vertical"} flexItem/>
                <InstanceMain tab={tab} database={activeInstance.database}/>
            </Box>
        )
    }
}
