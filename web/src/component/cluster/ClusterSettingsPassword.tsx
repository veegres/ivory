import {shortUuid} from "../../app/utils";
import {Autocomplete, Box, TextField} from "@mui/material";
import React, {useEffect, useMemo, useState} from "react";
import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query";
import {clusterApi, credentialApi} from "../../app/api";
import {Cluster, ClusterMap, CredentialMap, CredentialType} from "../../app/types";
import {useToast} from "../../app/hooks";

const keys = {
    [CredentialType.POSTGRES]: "postgresCredId",
    [CredentialType.PATRONI]: "patroniCredId",
}

type Value = {
    key: string
    short: string
    username?: string
}

type Props = {
    type: CredentialType
    cluster: Cluster
    credId: string
    label: string
}

export function ClusterSettingsPassword(props: Props) {
    const { type, label, cluster, credId } = props
    const passKey = keys[type]
    const [value, setValue] = useState<Value | null>(null)
    const [inputValue, setInputValue] = useState(credId);
    const { onError } = useToast()

    const query = useQuery(["credentials", type], () => credentialApi.list(type))
    const queryClient = useQueryClient();
    const updateCluster = useMutation(clusterApi.update, {
        onSuccess: (data) => {
            const map = queryClient.getQueryData<ClusterMap>(["cluster/list"]) ?? {} as ClusterMap
            map[data.name] = data
            queryClient.setQueryData<ClusterMap>(["cluster/list"], map)
        },
        onError,
    })

    const map = useMemo(() => query.data ?? {}, [query.data])
    const options = useMemo(() => handleMemoOptions(map), [map])
    const isPasswordRemoved = !!value && !map[credId]

    useEffect(handleEffectValue, [credId, map])

    return (
        <Autocomplete
            options={options}
            value={value}
            onChange={(_, value) => handleOnChange(value)}
            inputValue={inputValue}
            loading={query.isLoading || updateCluster.isLoading}
            onInputChange={(_, value) => setInputValue(value)}
            getOptionLabel={(option) => `${option.username ?? "***"} [${option.short}]`}
            isOptionEqualToValue={(option, value) => option.key === value.key}
            renderOption={(props, option) => <Box component={"li"} {...props}>{option.username} [{option.short}]</Box>}
            renderInput={(params) => (
                <TextField
                    {...params}
                    size={"small"}
                    label={label}
                    error={isPasswordRemoved}
                    helperText={isPasswordRemoved && "password was removed"}
                />
            )}
        />
    )

    function handleOnChange(value: Value | null) {
        updateCluster.mutate({...cluster, [passKey]: value?.key})
    }

    function handleEffectValue() {
        if (credId) setValue({ key: credId, short: shortUuid(credId), username: map[credId]?.username })
        else setValue(null)
    }

    function handleMemoOptions(map: CredentialMap) {
        return Object.entries(map).map(([key, value]) => ({ key, short: shortUuid(key), username: value.username }))
    }
}
