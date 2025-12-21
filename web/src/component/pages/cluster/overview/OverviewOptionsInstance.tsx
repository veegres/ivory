import {
    Autocomplete,
    AutocompleteChangeReason,
    Box,
    TextField,
    ToggleButton,
    ToggleButtonGroup,
    Tooltip
} from "@mui/material"
import {useMemo, useState} from "react"

import {Instance} from "../../../../api/cluster/type"
import {SxPropsMap} from "../../../../app/type"
import {getDomain} from "../../../../app/utils"
import {useStoreAction} from "../../../../provider/StoreProvider"

const SX: SxPropsMap = {
    box: {display: "flex", gap: 1},
    autocomplete: {flex: "auto"},
    toggle: {width: "30px"}
}

type Props = {
    instances: { [p: string]: Instance },
    mainInstance?: Instance,
    detectBy?: Instance,
}

export function OverviewOptionsInstance(props: Props) {
    const {setClusterDetection} = useStoreAction
    const {detectBy, instances, mainInstance} = props

    const value = getDomain(detectBy?.sidecar ?? mainInstance?.sidecar ?? {host: "-", port: 0})
    const [inputValue, setInputValue] = useState<string | undefined>(value)

    const options = useMemo(() => Object.keys(instances), [instances])

    return (
        <Box sx={SX.box}>
            <Autocomplete
                sx={SX.autocomplete}
                size={"small"}
                options={options}
                value={value}
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
                        selected={!detectBy}
                        onClick={() => setClusterDetection(undefined)}>
                        A
                    </ToggleButton>
                </Tooltip>
                <Tooltip title={"MANUAL"} placement={"top"}>
                    <ToggleButton
                        sx={SX.toggle}
                        value={"manual"}
                        selected={!!detectBy}
                        onClick={() => setClusterDetection(mainInstance)}>
                        M
                    </ToggleButton>
                </Tooltip>
            </ToggleButtonGroup>
        </Box>
    )

    function handleOnChange(value: string, reason: AutocompleteChangeReason) {
        if (value && reason === "selectOption") setClusterDetection(instances[value])
    }
}
