import {Autocomplete, Stack, TextField} from "@mui/material";
import React from "react";
import {CredentialType} from "../../app/types";
import {ClusterSettingsPassword} from "./ClusterSettingsPassword";
import {TabProps} from "./Cluster";
import {ClusterSettingsInstance} from "./ClusterSettingsInstance";

const SX = {
    settings: { width: "250px", gap: "12px", padding: "8px 0" }
}

export function ClusterSettings({info}: TabProps) {
    const { instances, instance, cluster, detection } = info

    const postgresCredId = cluster.postgresCredId ?? ""
    const patroniCredId = cluster.patroniCredId ?? ""

    return (
        <Stack sx={SX.settings}>
            <ClusterSettingsInstance instance={instance.api_domain} instances={instances} detection={detection} />
            <ClusterSettingsPassword label={"Postgres Password"} type={CredentialType.POSTGRES} credId={postgresCredId} cluster={cluster} />
            <ClusterSettingsPassword label={"Patroni Password"} type={CredentialType.PATRONI} credId={patroniCredId} cluster={cluster} />
            <Autocomplete
                value={cluster.certsId}
                options={[]}
                disabled={true}
                renderInput={(params) => <TextField {...params} size={"small"} label={"Patroni Certs (coming soon)"} />}
            />
            <Autocomplete
                multiple
                options={[]}
                disabled={true}
                renderInput={(params) => <TextField {...params} size={"small"} label={"Tags (coming soon)"} />}
            />
        </Stack>
    )
}
