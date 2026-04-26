import {
    Autocomplete, AutocompleteChangeReason, Box, TextField,
    ToggleButton, ToggleButtonGroup, Tooltip
} from "@mui/material"
import {useMemo, useState} from "react"

import {NodeOverview} from "../../../../api/cluster/type"
import {SxPropsMap} from "../../../../app/type"
import {useStoreAction} from "../../../../provider/StoreProvider"

const SX: SxPropsMap = {
    box: {display: "flex", gap: 1},
    autocomplete: {flex: "auto"},
    toggle: {width: "30px"}
}

type Props = {
    nodes: NodeOverview,
    mainKeeper?: string,
    manualKeeper?: string,
}

export function OverviewOptionsNode(props: Props) {
    const {setClusterDetection} = useStoreAction
    const {manualKeeper, nodes, mainKeeper} = props

    const value = manualKeeper ?? mainKeeper ?? "-"
    const [inputValue, setInputValue] = useState<string | undefined>(value)

    const options = useMemo(() => Object.keys(nodes), [nodes])

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
                renderInput={(params) => <TextField {...params} label={"Main Keeper"}/>}
            />
            <ToggleButtonGroup size={"small"}>
                <Tooltip title={"AUTO"} placement={"top"}>
                    <ToggleButton
                        sx={SX.toggle}
                        value={"auto"}
                        selected={!manualKeeper}
                        onClick={() => setClusterDetection(undefined)}>
                        A
                    </ToggleButton>
                </Tooltip>
                <Tooltip title={"MANUAL"} placement={"top"}>
                    <ToggleButton
                        sx={SX.toggle}
                        value={"manual"}
                        selected={!!manualKeeper}
                        onClick={() => setClusterDetection(mainKeeper ?? options[0])}>
                        M
                    </ToggleButton>
                </Tooltip>
            </ToggleButtonGroup>
        </Box>
    )

    function handleOnChange(value: string, reason: AutocompleteChangeReason) {
        if (value && reason === "selectOption") setClusterDetection(value)
    }
}
