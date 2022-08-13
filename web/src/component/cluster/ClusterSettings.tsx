import {Autocomplete, Stack, TextField} from "@mui/material";
import React from "react";
import {Cluster, CredentialType, Instance} from "../../app/types";
import {ClusterPassword} from "./ClusterPassword";

const SX = {
    settings: { width: "250px", gap: 1 }
}

type Props = {
    cluster: Cluster
    instance?: Instance
}

export function ClusterSettings(props: Props) {
    const { cluster, instance } = props

    const passPostgres = cluster.postgresCredId ?? ""
    const passPatroni = cluster.patroniCredId ?? ""

    return (
        <Stack sx={SX.settings}>
            <Autocomplete
                value={instance?.api_domain ?? ""}
                options={cluster.nodes}
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
