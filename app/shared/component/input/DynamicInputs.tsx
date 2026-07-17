import {Box, Chip, FormControl, OutlinedInput, SxProps, Theme} from "@mui/material"
import {ChangeEvent, memo, ReactNode, useCallback, useMemo} from "react"

import {ColorsMap, SxPropsMap} from "../../helper/HelperType"
import {SxPropsFormatter} from "../../helper/HelperUtils"

const SX: SxPropsMap = {
    box: {display: "grid", gap: 1},
    chip: {width: "100%", borderRadius: 2},
    input: {height: "32px"},
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
    helper?: ReactNode,
}

export const DynamicInputs = memo(function DynamicInputs(props: Props) {
    const {inputs, editable, placeholder, onChange, colors, InputProps, minLength = 1, InputSize = "263px", helper} = props
    const colorsMap = colors ?? EMPTY_COLORS

    const mergedSx = useMemo(() => SxPropsFormatter.merge(SX.input, InputProps), [InputProps])
    const nodesWithEmptyElement = useMemo(handleMemoEmptyNodes, [inputs, minLength])

    return (
        <Box sx={[SX.box, {gridTemplateColumns: `repeat(auto-fill, ${InputSize})`}]}>
            {editable ? (
                nodesWithEmptyElement.map((input, index) => (
                    <InputItem
                        key={index}
                        index={index}
                        value={input}
                        inputs={nodesWithEmptyElement}
                        placeholder={placeholder}
                        mergedSx={mergedSx}
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
    mergedSx: SxProps<Theme>,
    color: "success" | "primary" | "error" | "warning" | "info" | undefined,
    onChange: (values: string[]) => void,
    helper?: ReactNode,
}

const InputItem = memo(function InputItem(props: InputItemProps) {
    const {index, value, inputs, placeholder, mergedSx, color, helper, onChange} = props

    const handleChange = useCallback((event: ChangeEvent<HTMLInputElement>) => {
        const updated = [...inputs]
        updated[index] = event.target.value
        const lastSymbolIndex = updated.findLastIndex(s => s !== "")
        onChange(updated.slice(0, lastSymbolIndex + 1))
    }, [inputs, index, onChange])

    return (
        <FormControl color={color} focused={!!color}>
            <OutlinedInput
                sx={mergedSx}
                type={"string"}
                placeholder={`${placeholder}${index + 1}`}
                size={"small"}
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

