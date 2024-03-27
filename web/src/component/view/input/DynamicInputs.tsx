import {Box, Chip, FormControl, OutlinedInput, SxProps, Theme} from "@mui/material";
import {ColorsMap, SxPropsMap} from "../../../type/general";
import {SxPropsFormatter} from "../../../app/utils";

const SX: SxPropsMap = {
    chip: {width: '100%'},
    input: {height: '32px'},
    box: {display: "grid", gridTemplateColumns: "repeat(auto-fill, 175px)", gap: 1},
    label: {display: "flex", gap: 1},
}

type Props = {
    inputs: string[],
    colors?: ColorsMap,
    placeholder: string,
    editable: boolean,
    onChange: (values: string[]) => void,
    InputProps?: SxProps<Theme>,
}

export function DynamicInputs(props: Props) {
    const {inputs, editable, placeholder, onChange, colors, InputProps} = props
    const colorsMap = colors ?? {}
    return (
        <Box sx={SX.box}>
            {editable ? renderInputs() : renderChips()}
        </Box>
    )

    function renderInputs() {
        const nodesWithEmptyElement = inputs[inputs.length - 1] === '' ? inputs : [...inputs, '']
        return nodesWithEmptyElement.map((input, index) => {
            const color = colorsMap[input.toLowerCase()]
            return (
                <FormControl key={index} color={color} focused={!!color}>
                    <OutlinedInput
                        sx={SxPropsFormatter.merge(SX.input, InputProps)}
                        type={"string"}
                        placeholder={`${placeholder}${index+1}`}
                        size={"small"}
                        value={input}
                        onChange={(event) => handleChange(nodesWithEmptyElement, index, event.target.value)}
                    />
                </FormControl>
            )
        })
    }

    function renderChips() {
        return inputs.map((input, index) => (
            <Chip
                key={index}
                sx={SX.chip}
                color={colorsMap[input.toLowerCase()]}
                label={input}
                disabled={!input}
                variant={"outlined"}
            />
        ))
    }

    function handleChange(array: string[], index: number, value: string) {
        array[index] = value
        onChange(array.filter(node => !!node))
    }
}
