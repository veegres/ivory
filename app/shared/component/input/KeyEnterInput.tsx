import {VisibilityOffOutlined, VisibilityOutlined} from "@mui/icons-material"
import {FormControl, FormHelperText, IconButton, InputAdornment, InputLabel, OutlinedInput} from "@mui/material"
import {InputProps as StandardInputProps} from "@mui/material"
import {useState} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    margin: {margin: "0px"},
}

type Props = {
    label: string,
    hidden?: boolean,
    value?: string,
    required?: boolean,
    disabled?: boolean,
    error?: boolean,
    helperText?: string,
    onChange: StandardInputProps["onChange"],
    onEnterPress?: () => void,
}

export function KeyEnterInput(props: Props) {
    const {label, onChange, hidden = false, required = true, disabled = false, error = false} = props
    const {helperText, onEnterPress, value} = props
    const [show, setShow] = useState(false)

    return (
        <FormControl
            sx={SX.margin}
            fullWidth
            required={required}
            disabled={disabled}
            error={error}
            margin={"normal"}
        >
            <InputLabel>{label}</InputLabel>
            <OutlinedInput
                value={value}
                type={hidden && !show ? "password" : "text"}
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
                <IconButton size={"small"} onClick={() => setShow(!show)}>
                    {show ? <VisibilityOffOutlined fontSize={"small"}/> : <VisibilityOutlined fontSize={"small"}/>}
                </IconButton>
            </InputAdornment>
        )
    }

    function handleKeyPress(key: string) {
        if (key === "Enter" && onEnterPress) onEnterPress()
    }
}
