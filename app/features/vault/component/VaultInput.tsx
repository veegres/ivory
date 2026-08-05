import {FormControl, InputLabel, OutlinedInput, Theme} from "@mui/material"
import {SystemStyleObject} from "@mui/system"
import {forwardRef} from "react"

import {SxPropsMap} from "../../../shared/helper/HelperType"

const SX: SxPropsMap = {
    input: {
        height: "32px",
        "& input": {fontSize: "14px", fontFamily: "monospace"},
        "& .MuiOutlinedInput-notchedOutline": {borderColor: "divider"},
        "&.Mui-disabled .MuiOutlinedInput-notchedOutline": {borderColor: "divider"},
    },
}

type Props = {
    sx?: SystemStyleObject<Theme>,
    label: string,
    type: string,
    value: string,
    disabled: boolean,
    onChange: (value: string) => void
}

export const VaultInput = forwardRef<HTMLDivElement, Props>(function VaultInput(props, ref) {
    const {sx, label, type, value, disabled, onChange, ...rest} = props
    return (
        <FormControl ref={ref} sx={sx} disabled={disabled} {...rest}>
            <InputLabel shrink>{label}</InputLabel>
            <OutlinedInput
                sx={SX.input}
                notched
                label={label}
                type={type}
                value={value}
                onChange={(e) => onChange(e.target.value)}
            />
        </FormControl>
    )
})
