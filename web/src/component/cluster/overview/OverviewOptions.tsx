import {Divider, Stack} from "@mui/material";
import {TabProps} from "./Overview";
import {OverviewOptionsInstance} from "./OverviewOptionsInstance";
import {getDomain} from "../../../app/utils";
import {SxPropsMap} from "../../../type/common";
import {Options} from "../../shared/options/Options";

const SX: SxPropsMap = {
    settings: {width: "250px", gap: "12px", padding: "8px 0"}
}

export function OverviewOptions({info}: TabProps) {
    const {combinedInstanceMap, defaultInstance, cluster, detection} = info

    return (
        <Stack sx={SX.settings}>
            <OverviewOptionsInstance instance={getDomain(defaultInstance.sidecar)} instances={combinedInstanceMap} detection={detection}/>
            <Divider variant={"middle"}/>
            <Options cluster={cluster} />
        </Stack>
    )
}
