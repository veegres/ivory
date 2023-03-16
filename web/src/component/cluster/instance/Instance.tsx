import {Box} from "@mui/material";
import {useStore} from "../../../provider/StoreProvider";
import {InfoAlert} from "../../view/InfoAlert";
import {PageBlock} from "../../view/PageBlock";
import {SxPropsMap} from "../../../type/common";
import {useState} from "react";
import {InstanceMain} from "./InstanceMain";
import {InstanceInfo} from "./InstanceInfo";
import {InstanceTabType} from "../../../type/Instance";

const SX: SxPropsMap = {
    content: {display: "flex", gap: 3},
}

export function Instance() {
    const {store: {activeInstance}, isClusterOverviewOpen} = useStore()
    const [tab, setTab] = useState(InstanceTabType.CHART)

    return (
        <PageBlock withPadding visible={isClusterOverviewOpen()}>
            {renderContent()}
        </PageBlock>
    )

    function renderContent() {
        if (!activeInstance) return <InfoAlert text={"Please, select a instance to see the information!"}/>

        return (
            <Box sx={SX.content}>
                <InstanceInfo instance={activeInstance} tab={tab} onTab={setTab}/>
                <InstanceMain tab={tab} database={activeInstance.database}/>
            </Box>
        )
    }
}
