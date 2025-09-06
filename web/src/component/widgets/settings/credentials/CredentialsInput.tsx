import {FormControl, FormHelperText, InputLabel, OutlinedInput} from "@mui/material";
import {useEffect, useState} from "react";
import {SxPropsMap} from "../../../../api/management/type";

const SX: SxPropsMap = {
    input: {height: "36px"},
}

type Props = {
    label: string,
    type: string,
    value: string,
    disabled: boolean,
    error: boolean,
    onChange: (value: string) => void
}

export function CredentialsInput(props: Props) {
    const {label, type, value, disabled, onChange, error: initError = false} = props
    const [error, setError] = useState(false)

    useEffect(() => { if (!initError) setError(false) }, [initError])

    const isErrorShown = !disabled && error
    return (
        <FormControl disabled={disabled} error={isErrorShown}>
            <InputLabel shrink>{label}</InputLabel>
            <OutlinedInput
                sx={SX.input}
                notched
                label={label}
                type={type}
                value={value}
                onChange={(e) => handleOnChange(e.target.value)}
            />
            <FormHelperText>{isErrorShown ? "empty field" : " "}</FormHelperText>
        </FormControl>
    )

    function handleOnChange(value: string) {
        if (value) setError(false); else setError(true)
        onChange(value)
    }
}
