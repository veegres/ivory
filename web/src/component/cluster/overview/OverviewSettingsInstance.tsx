import {Autocomplete, Box, TextField, ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material";
import React, {useMemo, useState} from "react";
import {DetectionType, OverviewMap} from "../../../app/types";
import {useStore} from "../../../provider/StoreProvider";

const SX = {
    box: { display: "flex", gap: 1 },
    autocomplete: { flex: "auto" },
    toggle: { width: "30px" }
}

type Props = {
    instance: string
    instances: OverviewMap
    detection: DetectionType
}

export function OverviewSettingsInstance(props: Props) {
    const { setClusterInstance, setClusterDetection } = useStore()
    const { instance, instances, detection } = props
    const [inputValue, setInputValue] = useState<string | undefined>(instance);

    const options = useMemo(() => Object.keys(instances), [instances])

    return (
        <Box sx={SX.box}>
            <Autocomplete
                sx={SX.autocomplete}
                options={options}
                value={instance}
                disableClearable
                onChange={(_, value) => handleOnChange(value)}
                inputValue={inputValue}
                isOptionEqualToValue={(option, value) => option === value}
                onInputChange={(_, value) => setInputValue(value)}
                renderInput={(params) => <TextField {...params} size={"small"} label={"Default Instance"} />}
            />
            <ToggleButtonGroup size={"small"}>
                <Tooltip title={"AUTO"} placement={"top"}>
                    <ToggleButton
                        sx={SX.toggle}
                        value={"auto"}
                        selected={detection === "auto"}
                        onClick={() => setClusterDetection("auto")}>
                        A
                    </ToggleButton>
                </Tooltip>
                <Tooltip title={"MANUAL"} placement={"top"}>
                    <ToggleButton
                        sx={SX.toggle}
                        value={"manual"}
                        selected={detection === "manual"}
                        onClick={() => setClusterDetection("manual")}>
                        M
                    </ToggleButton>
                </Tooltip>
            </ToggleButtonGroup>
        </Box>
    )

    function handleOnChange(value?: string) {
        if (value) setClusterInstance(instances[value])
    }
}
