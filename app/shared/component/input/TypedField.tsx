import {TextField} from "@mui/material"

type Props = {
    label: string,
    example?: string,
    type: "text" | "port",
    value: string,
    disabled: boolean,
    onChange: (value: string) => void,
}

// TypedField renders one editable field; port fields take numbers and the
// example becomes helper text.
export function TypedField(props: Props) {
    const {label, example, type, value, disabled, onChange} = props
    return (
        <TextField
            fullWidth={true}
            size={"small"}
            type={type === "port" ? "number" : "text"}
            disabled={disabled}
            label={label}
            helperText={example ? `Example: ${example}` : undefined}
            value={value}
            onChange={(e) => onChange(e.target.value)}
        />
    )
}
