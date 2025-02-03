import {Menu} from "../../settings/menu/Menu";
import {Box, Stack} from "@mui/material";
import {List as ClusterList} from "./list/List";
import {Overview as ClusterOverview} from "./overview/Overview";
import {Instance as ClusterInstance} from "./instance/Instance";
import {SxPropsMap} from "../../../type/general";

const SX: SxPropsMap = {
    stack: {width: "100%", height: "100%", gap: 4}
}

export function ClusterBody() {
    return (
        <Box>
            <Menu/>
            <Stack sx={SX.stack}>
                <ClusterList/>
                <ClusterOverview/>
                <ClusterInstance/>
            </Stack>
        </Box>
    )
}
