import {shortUuid} from "../../app/utils";
import {Autocomplete, Box, TextField} from "@mui/material";
import React, {useMemo, useState} from "react";
import {useQuery, UseQueryResult} from "react-query";
import {credentialApi} from "../../app/api";
import {CredentialMap, CredentialType} from "../../app/types";

type Value = {
    key: string
    short: string
    username: string
}

type Props = {
    type: CredentialType
    pass: string
    label: string
}

export function ClusterSettingsPassword(props: Props) {
    const { type, pass, label } = props
    const [value, setValue] = useState<Value | null>(null);
    const [inputValue, setInputValue] = useState(pass);

    const query = useQuery(["credentials", type], () => credentialApi.get(type))
    const options = useMemo(() => handleMemoOptions(query), [query])

    return (
        <Autocomplete
            options={options}
            value={value}
            onChange={(_, value) => setValue(value)}
            inputValue={inputValue}
            onInputChange={(_, value) => setInputValue(value)}
            getOptionLabel={(option) => `${option.username} [${option.short}]`}
            isOptionEqualToValue={(option, value) => option.key === value.key}
            renderOption={(props, option) => <Box component="li" {...props}>{option.username} [{option.short}]</Box>}
            renderInput={(params) => <TextField {...params} size={"small"} label={label} />}
        />
    )

    function handleMemoOptions(query: UseQueryResult<CredentialMap>) {
        return Object.entries(query.data ?? {}).map(([key, value]) => ({ key, short: shortUuid(key), username: value.username }))
    }
}
