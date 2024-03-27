import {
    Autocomplete,
    AutocompleteChangeReason,
    Box,
    TextField,
    ToggleButton,
    ToggleButtonGroup,
    Tooltip
} from "@mui/material";
import {useMemo, useState} from "react";
import {useStoreAction} from "../../../../provider/StoreProvider";
import {SxPropsMap} from "../../../../type/general";
import {InstanceMap} from "../../../../type/instance";
import {DetectionType} from "../../../../type/cluster";

const SX: SxPropsMap = {
    box: {display: "flex", gap: 1},
    autocomplete: {flex: "auto"},
    toggle: {width: "30px"}
}

type Props = {
    instance: string
    instances: InstanceMap
    detection: DetectionType
}

export function OverviewOptionsInstance(props: Props) {
    const {setClusterInstance, setClusterDetection} = useStoreAction()
    const {instance, instances, detection} = props
    const [inputValue, setInputValue] = useState<string | undefined>(instance);

    const options = useMemo(() => Object.keys(instances), [instances])

    return (
        <Box sx={SX.box}>
            <Autocomplete
                sx={SX.autocomplete}
                size={"small"}
                options={options}
                value={instance}
                disableClearable
                onChange={(_, value, reason) => handleOnChange(value, reason)}
                inputValue={inputValue}
                isOptionEqualToValue={(option, value) => option === value}
                onInputChange={(_, value) => setInputValue(value)}
                renderInput={(params) => <TextField {...params} label={"Main Instance"}/>}
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

    function handleOnChange(value: string, reason: AutocompleteChangeReason) {
        if (value && reason === "selectOption") setClusterInstance(instances[value])
    }
}
