import {Autocomplete, Stack, TextField} from "@mui/material";
import React from "react";
import {CredentialType} from "../../app/types";
import {ClusterPassword} from "./ClusterPassword";
import {TabProps} from "./Cluster";

const SX = {
    settings: { width: "250px", gap: 1 }
}

export function ClusterSettings({info}: TabProps) {
    const { instances, instance, cluster } = info

    const passPostgres = cluster.postgresCredId ?? ""
    const passPatroni = cluster.patroniCredId ?? ""

    return (
        <Stack sx={SX.settings}>
            <Autocomplete
                value={instance?.api_domain ?? ""}
                options={Object.keys(instances)}
                renderInput={(params) => <TextField {...params} size={"small"} label={"Default Instance"} />}
            />
            <ClusterPassword label={"Postgres Password"} type={CredentialType.POSTGRES} pass={passPostgres} />
            <ClusterPassword label={"Patroni Password"} type={CredentialType.PATRONI} pass={passPatroni} />
            <Autocomplete
                value={cluster.certsId}
                options={[]}
                disabled={true}
                renderInput={(params) => <TextField {...params} size={"small"} label={"Certs (not implemented)"} />}
            />
            <Autocomplete
                multiple
                options={[]}
                disabled={true}
                renderInput={(params) => <TextField {...params} size={"small"} label={"Tags (not implemented)"} />}
            />
        </Stack>
    )
}
