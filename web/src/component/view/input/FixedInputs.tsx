import {Box, FormControl, OutlinedInput} from "@mui/material";
import {SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    input: {height: '32px'},
    box: {display: "grid", gridTemplateColumns: "repeat(auto-fill, 150px)", gap: 1},
}

type Props = {
    inputs: string[],
    placeholders: string[],
    onChange: (values: string[]) => void,
}

export function FixedInputs(props: Props) {
    const {inputs, placeholders, onChange} = props

    if (inputs.length != placeholders.length) {
        throw Error("input length should be equal to placeholder length")
    }

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
                        onChange={(event) => handleChange(index, event.target.value)}
                    />
                </FormControl>
            ))}
        </Box>
    )

    function handleChange(index: number, value: string) {
        inputs[index] = value
        onChange([...inputs])
    }
}
