import {shortUuid} from "../../app/utils";
import {Autocomplete, Box, TextField} from "@mui/material";
import React, {useState} from "react";
import {useQuery} from "react-query";
import {credentialApi} from "../../app/api";
import {CredentialType} from "../../app/types";

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

export function ClusterPassword(props: Props) {
    const { type, pass, label } = props
    const [value, setValue] = useState<Value | null>(null);
    const [inputValue, setInputValue] = useState(pass);

    const query = useQuery(["credentials", type], () => credentialApi.get(type))
    const list = Object.entries(query.data ?? {}).map(([key, value]) => ({ key, short: shortUuid(key), username: value.username }))


    return (
        <Autocomplete
            options={list}
            value={value}
            onChange={(event, value) => setValue(value)}
            inputValue={inputValue}
            onInputChange={(event, value) => setInputValue(value)}
            getOptionLabel={(option) => `${option.username} [${option.short}]`}
            isOptionEqualToValue={(option, value) => option.key === value.key}
            renderOption={(props, option) => <Box component="li" {...props}>{option.username} [{option.short}]</Box>}
            renderInput={(params) => <TextField {...params} size={"small"} label={label} />}
        />
    )
}
