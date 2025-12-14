import {Stack} from "@mui/material"

import {Permission} from "../../../api/permission/type"
import {SxPropsMap} from "../../../app/type"
import {Access} from "../../widgets/access/Access"
import {Menu} from "../../widgets/settings/menu/Menu"
import {Instance as ClusterInstance} from "./instance/Instance"
import {List as ClusterList} from "./list/List"
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
                <Access permission={Permission.ViewInstanceOverview}>
                    <ClusterOverview/>
                    <ClusterInstance/>
                </Access>
            </Stack>
        </>
    )
}
