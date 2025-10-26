import {Visibility, VisibilityOff} from "@mui/icons-material"
import {FormControl, FormHelperText, IconButton, InputAdornment, InputLabel, OutlinedInput} from "@mui/material"
import {InputProps as StandardInputProps} from "@mui/material/Input/Input"
import {useState} from "react"

import {SxPropsMap} from "../../../app/type"

const SX: SxPropsMap = {
    margin: {margin: "0px"},
}

type Props = {
    label: string,
    hidden?: boolean,
    value?: string,
    required?: boolean,
    disabled?: boolean,
    helperText?: string,
    onChange: StandardInputProps["onChange"],
    onEnterPress?: () => void,
}

export function KeyEnterInput(props: Props) {
    const {label, onChange, hidden = false, required = true, disabled = false, helperText, onEnterPress, value} = props
    const [showPassword, setShowPassword] = useState(false)

    return (
        <FormControl sx={SX.margin} fullWidth required={required} disabled={disabled} size={"small"} margin={"normal"}>
            <InputLabel>{label}</InputLabel>
            <OutlinedInput
                value={value}
                type={hidden && !showPassword ? "password" : "text"}
                endAdornment={renderAdornment()}
                label={label}
                autoComplete={"new-password"}
                onChange={onChange}
                onKeyUp={(e) => handleKeyPress(e.key)}
            />
            {helperText && <FormHelperText>{helperText}</FormHelperText>}
        </FormControl>
    )

    function renderAdornment() {
        if (!hidden) return

        return (
            <InputAdornment position={"end"}>
                <IconButton onClick={() => setShowPassword(!showPassword)}>
                    {showPassword ? <VisibilityOff /> : <Visibility />}
                </IconButton>
            </InputAdornment>
        )
    }

    function handleKeyPress(key: string) {
        if (key === "Enter" && onEnterPress) onEnterPress()
    }
}
