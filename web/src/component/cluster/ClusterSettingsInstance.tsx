import {Autocomplete, TextField} from "@mui/material";
import React, {useMemo, useState} from "react";
import {InstanceMap} from "../../app/types";
import {useStore} from "../../provider/StoreProvider";

type Props = {
    instance: string
    instances: InstanceMap
}

export function ClusterSettingsInstance(props: Props) {
    const { setClusterInstance } = useStore()
    const { instance, instances } = props
    const [inputValue, setInputValue] = useState<string | undefined>(instance);

    const options = useMemo(() => Object.keys(instances), [instances])

    return (
        <Autocomplete
            options={options}
            value={instance}
            disableClearable
            onChange={(_, value) => handleOnChange(value)}
            inputValue={inputValue}
            isOptionEqualToValue={(option, value) => option === value}
            onInputChange={(_, value) => setInputValue(value)}
            renderInput={(params) => <TextField {...params} size={"small"} label={"Default Instance"} />}
        />
    )

    function handleOnChange(value?: string) {
        if (value) setClusterInstance(instances[value])
    }
}
