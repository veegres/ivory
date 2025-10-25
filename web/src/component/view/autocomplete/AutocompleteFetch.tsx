import {Autocomplete, AutocompleteRenderInputParams, TextField} from "@mui/material";
import {UseQueryResult} from "@tanstack/react-query";
import {useState} from "react";

import {QueryConnection} from "../../../api/query/type";
import {useDebounce} from "../../../hook/Debounce";

type Props = {
    onUpdate: (option: string | null) => void,
    disabled?: boolean,
    label?: string,
    placeholder?: string,
    variant?: "standard" | "filled" | "outlined",
    margin?: "dense" | "normal" | "none",
    size?: "small" | "medium",
    padding?: string,
    value?: string | null,
    connection: QueryConnection,
    params?: any,
    useFetch: (connection: QueryConnection, params: any, enabled?: boolean) => UseQueryResult<string[], Error>
}

export function AutocompleteFetch(props: Props) {
    const {label, placeholder, variant, margin, padding, size} = props
    const {onUpdate, useFetch, disabled = false, value, connection, params} = props

    const [inputValue, setInputValue] = useState("")
    const debInputValue = useDebounce(inputValue)

    const query = useFetch(connection, {name: debInputValue, ...params}, !disabled)
    const options = query.data ?? []

    return (
        <Autocomplete
            fullWidth
            freeSolo
            size={size}
            renderInput={renderInput}
            value={value}
            onChange={(_, v) => onUpdate(v)}
            inputValue={inputValue}
            onInputChange={(_, v) => setInputValue(v)}
            disabled={disabled}
            options={options}
            loading={query.isPending}
        />
    )

    function renderInput(params: AutocompleteRenderInputParams) {
        const inputProps = padding ? {...params.InputProps, style: {padding}} : params.InputProps
        return (
            <TextField
                {...params}
                size={"small"}
                label={label}
                variant={variant}
                margin={margin}
                placeholder={placeholder}
                slotProps={{input: inputProps}}
            />
        )
    }
}
