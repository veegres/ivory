import {Stack} from "@mui/material";

import {SxPropsMap} from "../../../app/type";
import {Menu} from "../../widgets/settings/menu/Menu";
import {Instance as ClusterInstance} from "./instance/Instance";
import {List as ClusterList} from "./list/List";
import {Overview as ClusterOverview} from "./overview/Overview";

const SX: SxPropsMap = {
    stack: {width: "100%", height: "100%", gap: 4}
}

export function ClusterBody() {
    return (
        <>
            <Menu/>
            <Stack sx={SX.stack}>
                <ClusterList/>
                <ClusterOverview/>
                <ClusterInstance/>
            </Stack>
        </>
    )
}
