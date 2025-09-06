import {useState} from "react";
import {FormControl, IconButton, InputAdornment, InputLabel, OutlinedInput} from "@mui/material";
import {InputProps as StandardInputProps} from "@mui/material/Input/Input";
import {Visibility, VisibilityOff} from "@mui/icons-material";
import {SxPropsMap} from "../../../api/management/type";

const SX: SxPropsMap = {
    margin: {margin: "0px"},
}

type Props = {
    label: string,
    hidden?: boolean,
    value?: unknown,
    onChange: StandardInputProps['onChange'],
    onEnterPress?: () => void,
}

export function KeyEnterInput(props: Props) {
    const {label, onChange, hidden = false, onEnterPress, value} = props
    const [showPassword, setShowPassword] = useState(false)

    return (
        <FormControl sx={SX.margin} fullWidth required size={"small"} margin={"normal"}>
            <InputLabel>{label}</InputLabel>
            <OutlinedInput
                value={value}
                type={hidden && !showPassword ? "password" : "text"}
                endAdornment={renderAdornment()}
                label={label}
                autoComplete={"off"}
                onChange={onChange}
                onKeyUp={(e) => handleKeyPress(e.key)}
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

    function handleKeyPress(key: string) {
        if (key === 'Enter' && onEnterPress) onEnterPress()
    }
}
