import {Autocomplete, Divider, Stack, TextField} from "@mui/material";
import React from "react";
import {CertType, CredentialType, SxPropsMap} from "../../../app/types";
import {OverviewSettingsPassword} from "./OverviewSettingsPassword";
import {TabProps} from "./Overview";
import {OverviewSettingsInstance} from "./OverviewSettingsInstance";
import {OverviewSettingsCert} from "./OverviewSettingsCert";
import {getDomain} from "../../../app/utils";

const SX: SxPropsMap = {
    settings: { width: "250px", gap: "12px", padding: "8px 0" }
}

export function OverviewSettings({info}: TabProps) {
    const { instances, instance, cluster, detection } = info

    return (
        <Stack sx={SX.settings}>
            <OverviewSettingsInstance instance={getDomain(instance.sidecar)} instances={instances} detection={detection} />
            <Divider variant={"middle"}/>
            <OverviewSettingsPassword type={CredentialType.POSTGRES} cluster={cluster} />
            <OverviewSettingsPassword type={CredentialType.PATRONI} cluster={cluster} />
            <Divider variant={"middle"}/>
            <OverviewSettingsCert type={CertType.CLIENT_CA} cluster={cluster} />
            <OverviewSettingsCert type={CertType.CLIENT_CERT} cluster={cluster} />
            <OverviewSettingsCert type={CertType.CLIENT_KEY} cluster={cluster} />
            <Divider variant={"middle"}/>
            <Autocomplete
                multiple
                options={[]}
                disabled={true}
                renderInput={(params) => <TextField {...params} size={"small"} label={"Tags (coming soon)"} />}
            />
        </Stack>
    )
}
