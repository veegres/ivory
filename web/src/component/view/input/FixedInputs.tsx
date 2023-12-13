import {Box, FormControl, OutlinedInput} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {useEffect, useState} from "react";

const SX: SxPropsMap = {
    input: {height: '32px'},
    box: {display: "grid", gridTemplateColumns: "repeat(auto-fill, 150px)", gap: 1},
}

type Props = {
    placeholders: string[],
    onChange: (values: string[]) => void,
}

export function FixedInputs(props: Props) {
    const {placeholders, onChange} = props
    const [inputs, setInputs] = useState(placeholders.map(() => ''))

    useEffect(handleEffectPlaceholders, [placeholders]);
    useEffect(handleEffectInputs, [inputs, onChange]);

    return (
        <Box sx={SX.box}>
            {inputs.map((input, index) => (
                <FormControl key={index}>
                    <OutlinedInput
                        sx={SX.input}
                        type={"string"}
                        placeholder={placeholders[index]}
                        size={"small"}
                        value={input}
                        onChange={(event) => handleChange(inputs, index, event.target.value)}
                    />
                </FormControl>
            ))}
        </Box>
    )

    function handleEffectPlaceholders() {
        setInputs(placeholders.map(() => ''))
    }

    function handleEffectInputs() {
        onChange(inputs)
    }

    function handleChange(array: string[], index: number, value: string) {
        array[index] = value
        onChange(array.filter(node => !!node))
    }
}
