import {Box} from "@mui/material"

import {ClusterApi} from "../../../../features/cluster/api/ClusterRouter"
import {Refresher} from "../../../widgets/browser/Refresher"

// NOTE: the Box absorbs the size prop cloned in by ActionsLoader so the
// Refresher keeps its own default size; spacing comes from the toolbar gap
export function ListTableRefresher() {
    return (
        <Box>
            <Refresher queryKeys={[ClusterApi.list.keyCommon(), ClusterApi.overview.keyCommon()]}/>
        </Box>
    )
}