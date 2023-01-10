import {Autocomplete, Stack, TextField} from "@mui/material";
import React from "react";
import {CredentialType} from "../../../app/types";
import {OverviewSettingsPassword} from "./OverviewSettingsPassword";
import {TabProps} from "./Overview";
import {OverviewSettingsInstance} from "./OverviewSettingsInstance";
import {OverviewSettingsCert} from "./OverviewSettingsCert";
import {getDomain} from "../../../app/utils";

const SX = {
    settings: { width: "250px", gap: "12px", padding: "8px 0" }
}

export function OverviewSettings({info}: TabProps) {
    const { instances, instance, cluster, detection } = info

    const postgresCredId = cluster.postgresCredId ?? ""
    const patroniCredId = cluster.patroniCredId ?? ""
    const patroniCertId = cluster.certId ?? ""

    return (
        <Stack sx={SX.settings}>
            <OverviewSettingsInstance instance={getDomain(instance.sidecar)} instances={instances} detection={detection} />
            <OverviewSettingsPassword label={"Postgres Password"} type={CredentialType.POSTGRES} credId={postgresCredId} cluster={cluster} />
            <OverviewSettingsPassword label={"Patroni Password"} type={CredentialType.PATRONI} credId={patroniCredId} cluster={cluster} />
            <OverviewSettingsCert certId={patroniCertId} cluster={cluster} />
            <Autocomplete
                multiple
                options={[]}
                disabled={true}
                renderInput={(params) => <TextField {...params} size={"small"} label={"Tags (coming soon)"} />}
            />
        </Stack>
    )
}
