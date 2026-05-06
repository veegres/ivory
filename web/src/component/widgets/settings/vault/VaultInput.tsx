import {FormControl, InputLabel, OutlinedInput, Theme} from "@mui/material"
import {SystemStyleObject} from "@mui/system/styleFunctionSx/styleFunctionSx"

import {SxPropsMap} from "../../../../app/type"

const SX: SxPropsMap = {
    input: {height: "36px"},
}

type Props = {
    sx?: SystemStyleObject<Theme>,
    label: string,
    type: string,
    value: string,
    disabled: boolean,
    onChange: (value: string) => void
}

export function VaultInput(props: Props) {
    const {sx, label, type, value, disabled, onChange} = props
    return (
        <FormControl sx={sx} disabled={disabled}>
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
}
