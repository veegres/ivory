import {Box, FormControl, OutlinedInput} from "@mui/material";
import {useEffect} from "react";
import {SxPropsMap} from "../../../app/type";

const SX: SxPropsMap = {
    input: {height: '32px'},
    box: {display: "grid", gridTemplateColumns: "repeat(auto-fill, 150px)", gap: 1},
}

type Props = {
    placeholders: string[],
    values: string[],
    onChange: (values: string[]) => void,
}

export function FixedInputs(props: Props) {
    const {placeholders, values, onChange} = props

    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(handleEffectPlaceholders, [JSON.stringify(placeholders)]);

    return (
        <Box sx={SX.box}>
            {values.map((input, index) => (
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

    function handleEffectPlaceholders() {
        onChange(placeholders.map(() => ""))
    }

    function handleChange(index: number, value: string) {
        const tmp= [...values]
        tmp[index] = value
        onChange(tmp)
    }
}
