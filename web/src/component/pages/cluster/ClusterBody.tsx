import {Stack} from "@mui/material"

import {Feature} from "../../../api/feature"
import {SxPropsMap} from "../../../app/type"
import {Access} from "../../widgets/access/Access"
import {Menu} from "../../widgets/settings/menu/Menu"
import {List as ClusterList} from "./list/List"
import {Node as ClusterNode} from "./node/Node"
import {Overview as ClusterOverview} from "./overview/Overview"

const SX: SxPropsMap = {
    stack: {width: "100%", height: "100%", gap: 4},
}

export function ClusterBody() {
    return (
        <>
            <Menu/>
            <Stack sx={SX.stack}>
                <ClusterList/>
                <Access feature={Feature.ViewNodeDbOverview}>
                    <ClusterOverview/>
                    <ClusterNode/>
                </Access>
            </Stack>
        </>
    )
}
