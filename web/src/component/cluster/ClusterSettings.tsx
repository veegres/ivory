import {Autocomplete, Box, Stack, TextField} from "@mui/material";
import React from "react";
import {useStore} from "../../provider/StoreProvider";

const SX = {
    settings: { width: "250px", gap: 1 }
}

export function ClusterSettings() {
    const {store: {activeCluster, activeNode}, setStore} = useStore()

    console.log(activeNode);

    return (
        <Stack sx={SX.settings}>
            <Box>{activeCluster.name ?? "No Leader"}</Box>
            <Box>{activeNode ?? "No Node"}</Box>
            <Autocomplete
                options={[]}
                renderInput={(params) => <TextField {...params} size={"small"} label={"Default Node"} />}
            />
            <Autocomplete
                options={[]}
                renderInput={(params) => <TextField {...params} size={"small"} label={"Postgres Password"} />}
            />
            <Autocomplete
                options={[]}
                renderInput={(params) => <TextField {...params} size={"small"} label={"Patroni Password"} />}
            />
            <Autocomplete
                options={[]}
                renderInput={(params) => <TextField {...params} size={"small"} label={"Certs"} />}
            />
            <Autocomplete
                multiple
                options={[]}
                renderInput={(params) => <TextField {...params} size={"small"} label={"Tags"} />}
            />
        </Stack>
    )
}
