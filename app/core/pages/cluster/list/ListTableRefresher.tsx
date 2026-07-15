import {Box} from "@mui/material"

import {ClusterApi} from "../../../../features/cluster/api/ClusterRouter"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {Refresher} from "../../../widgets/browser/Refresher"

const SX: SxPropsMap = {
    box: {padding: "0px 5px"},
}

export function ListTableRefresher() {
    return (
        <Box sx={SX.box}>
            <Refresher queryKeys={[ClusterApi.list.keyCommon(), ClusterApi.overview.keyCommon()]}/>
        </Box>
    )
}