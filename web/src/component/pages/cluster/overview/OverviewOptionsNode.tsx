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

import {Node, NodeOverview} from "../../../../api/cluster/type"
import {SxPropsMap} from "../../../../app/type"
import {getDomain, initialNode} from "../../../../app/utils"
import {useStoreAction} from "../../../../provider/StoreProvider"

const SX: SxPropsMap = {
    box: {display: "flex", gap: 1},
    autocomplete: {flex: "auto"},
    toggle: {width: "30px"}
}

type Props = {
    nodes: NodeOverview,
    mainNode?: Node,
    detectBy?: Node,
}

export function OverviewOptionsNode(props: Props) {
    const {setClusterDetection} = useStoreAction
    const {detectBy, nodes, mainNode} = props

    const connection = detectBy?.connection ?? mainNode?.connection ?? {host: "-", keeperPort: 0}
    const value = getDomain(connection)
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
                        onClick={() => setClusterDetection(mainNode)}>
                        M
                    </ToggleButton>
                </Tooltip>
            </ToggleButtonGroup>
        </Box>
    )

    function handleOnChange(value: string, reason: AutocompleteChangeReason) {
        if (value && reason === "selectOption") setClusterDetection(nodes[value] ?? initialNode(value))
    }
}
