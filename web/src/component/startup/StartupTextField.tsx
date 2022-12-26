import {FormControl, IconButton, InputAdornment, InputLabel, OutlinedInput} from "@mui/material";
import {InputProps as StandardInputProps} from "@mui/material/Input/Input";
import {useState} from "react";
import {Visibility, VisibilityOff} from "@mui/icons-material";

type Props = {
    label: string,
    hidden?: boolean,
    onChange: StandardInputProps['onChange']
}

export function StartupTextField(props: Props) {
    const { label, onChange, hidden = false } = props
    const [showPassword, setShowPassword] = useState(false)
    return (
        <FormControl fullWidth required size={"small"} margin={"normal"}>
            <InputLabel>{label}</InputLabel>
            <OutlinedInput
                type={hidden && !showPassword ? "password" : "text"}
                endAdornment={renderAdornment()}
                label={label}
                autoComplete={"off"}
                onChange={onChange}
            />
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
}
