import {Menu} from "../../widgets/settings/menu/Menu";
import {Stack} from "@mui/material";
import {List as ClusterList} from "./list/List";
import {Overview as ClusterOverview} from "./overview/Overview";
import {Instance as ClusterInstance} from "./instance/Instance";
import {SxPropsMap} from "../../../api/management/type";

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
