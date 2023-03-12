import {Divider, Stack} from "@mui/material";
import {OverviewSettingsPassword} from "./OverviewSettingsPassword";
import {TabProps} from "./Overview";
import {OverviewSettingsInstance} from "./OverviewSettingsInstance";
import {OverviewSettingsCert} from "./OverviewSettingsCert";
import {getDomain} from "../../../app/utils";
import {OverviewSettingsTags} from "./OverviewSettingsTags";
import {SxPropsMap} from "../../../type/common";
import {PasswordType} from "../../../type/password";
import {CertType} from "../../../type/cert";

const SX: SxPropsMap = {
    settings: {width: "250px", gap: "12px", padding: "8px 0"}
}

export function OverviewSettings({info}: TabProps) {
    const {combinedInstanceMap, defaultInstance, cluster, detection} = info

    return (
        <Stack sx={SX.settings}>
            <OverviewSettingsInstance instance={getDomain(defaultInstance.sidecar)} instances={combinedInstanceMap} detection={detection}/>
            <Divider variant={"middle"}/>
            <OverviewSettingsPassword type={PasswordType.POSTGRES} cluster={cluster}/>
            <OverviewSettingsPassword type={PasswordType.PATRONI} cluster={cluster}/>
            <Divider variant={"middle"}/>
            <OverviewSettingsCert type={CertType.CLIENT_CA} cluster={cluster}/>
            <OverviewSettingsCert type={CertType.CLIENT_CERT} cluster={cluster}/>
            <OverviewSettingsCert type={CertType.CLIENT_KEY} cluster={cluster}/>
            <Divider variant={"middle"}/>
            <OverviewSettingsTags cluster={cluster}/>
        </Stack>
    )
}
