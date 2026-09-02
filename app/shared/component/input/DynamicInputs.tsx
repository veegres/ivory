import {Box, Chip, FormControl, OutlinedInput, SxProps, Theme} from "@mui/material"
import {ChangeEvent, memo, ReactNode, useCallback, useMemo} from "react"

import {ColorsMap, SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    box: {display: "grid", gap: 1},
    chip: {width: "100%", borderRadius: 2},
    label: {display: "flex", gap: 1},
    helper: {margin: "2px 14px 0px"},
}

const EMPTY_COLORS: ColorsMap = {}

type Props = {
    inputs: string[],
    colors?: ColorsMap,
    placeholder: string,
    editable: boolean,
    minLength?: number,
    onChange: (values: string[]) => void,
    InputSize?: string,
    InputProps?: SxProps<Theme>,
    size?: "small" | "medium",
    helper?: ReactNode,
}

export const DynamicInputs = memo(function DynamicInputs(props: Props) {
    const {inputs, editable, placeholder, onChange, colors, InputProps, size, minLength = 1, InputSize = "var(--size-input)", helper} = props
    const colorsMap = colors ?? EMPTY_COLORS

    const nodesWithEmptyElement = useMemo(handleMemoEmptyNodes, [inputs, minLength])

    return (
        <Box sx={[SX.box, {gridTemplateColumns: `repeat(auto-fill, minmax(min(${InputSize}, 100%), 1fr))`}]}>
            {editable ? (
                nodesWithEmptyElement.map((input, index) => (
                    <InputItem
                        key={index}
                        index={index}
                        value={input}
                        inputs={nodesWithEmptyElement}
                        placeholder={placeholder}
                        inputSx={InputProps}
                        size={size}
                        color={colorsMap[input.toLowerCase()]}
                        onChange={onChange}
                        helper={helper}
                    />
                ))
            ) : (
                inputs.map((input, index) => (
                    <ChipItem
                        key={index}
                        input={input}
                        color={colorsMap[input.toLowerCase()]}
                    />
                ))
            )}
        </Box>
    )

    function handleMemoEmptyNodes(){
        const tmp = inputs.length > minLength ? inputs : [...inputs, ...Array(minLength - inputs.length).fill("")]
        return tmp[tmp.length - 1] === "" ? tmp : [...tmp, ""]
    }
})

type InputItemProps = {
    index: number,
    value: string,
    inputs: string[],
    placeholder: string,
    inputSx?: SxProps<Theme>,
    size?: "small" | "medium",
    color: "success" | "primary" | "error" | "warning" | "info" | undefined,
    onChange: (values: string[]) => void,
    helper?: ReactNode,
}

const InputItem = memo(function InputItem(props: InputItemProps) {
    const {index, value, inputs, placeholder, inputSx, size, color, helper, onChange} = props

    const handleChange = useCallback((event: ChangeEvent<HTMLInputElement>) => {
        const updated = [...inputs]
        updated[index] = event.target.value
        const lastSymbolIndex = updated.findLastIndex(s => s !== "")
        onChange(updated.slice(0, lastSymbolIndex + 1))
    }, [inputs, index, onChange])

    return (
        <FormControl color={color} focused={!!color}>
            <OutlinedInput
                sx={inputSx}
                size={size}
                type={"string"}
                placeholder={`${placeholder}${index + 1}`}
                value={value}
                onChange={handleChange}
            />
            {helper && <Box sx={SX.helper}>{helper}</Box>}
        </FormControl>
    )
})

type ChipItemProps = {
    input: string,
    color: "success" | "primary" | "error" | "warning" | "info" | undefined,
}

const ChipItem = memo(function ChipItem(props: ChipItemProps) {
    const {input, color} = props
    return (
        <Chip
            sx={SX.chip}
            color={color}
            label={input}
            disabled={!input}
            variant={"outlined"}
        />
    )
})

