import {Box, Chip, FormControl, OutlinedInput} from "@mui/material";
import {ColorsMap} from "../../app/types";

const SX = {
    gridTemplateColumns: `repeat(auto-fill, minmax(150px, 1fr))`,
    chip: {width: '100%'},
    input: {height: '32px'},
}

type Props = {
    inputs: string[],
    colors?: ColorsMap,
    placeholder: string,
    editable: boolean,
    onChange: (nodes: string[]) => void,
}

export function DynamicInputs({ inputs, editable, placeholder, onChange, colors} : Props) {
    const colorsMap = colors ?? {}
    return (
        <Box display="grid" gridTemplateColumns={SX.gridTemplateColumns} gap={1}>
            {editable ? renderInputs(inputs) : renderChips(inputs)}
        </Box>
    )

    function renderInputs(nodes: string[]) {
        const nodesWithEmptyElement = nodes[nodes.length - 1] === '' ? nodes : [...nodes, '']

        return nodesWithEmptyElement.map((node, index) => (
            <FormControl key={index} color={colorsMap[node]} focused={!!colorsMap[node]}>
                <OutlinedInput
                    sx={SX.input}
                    type="string"
                    placeholder={`${placeholder} ${index}`}
                    size="small"
                    value={node}
                    onChange={(event) => handleChange(index, event.target.value)}
                />
            </FormControl>
        ))
    }

    function renderChips(nodes: string[]) {
        return nodes.map((node, index) => (
            <Chip
                key={index}
                sx={SX.chip}
                color={colorsMap[node]}
                label={node ? node : `Node ${index}`}
                disabled={!node}
                variant="outlined"
            />
        ))
    }

    function handleChange(index: number, value: string) {
        inputs[index] = value
        onChange(inputs.filter(node => !!node))
    }
}
