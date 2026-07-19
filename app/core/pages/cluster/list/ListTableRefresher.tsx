import {Box} from "@mui/material"

import {ClusterApi} from "../../../../features/cluster/api/ClusterRouter"
import {Refresher} from "../../../widgets/browser/Refresher"

export function ListTableRefresher() {
    return (
        <Box>
            <Refresher queryKeys={[ClusterApi.list.keyCommon(), ClusterApi.overview.keyCommon()]}/>
        </Box>
    )
}