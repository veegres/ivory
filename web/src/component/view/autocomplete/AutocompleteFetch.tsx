import {Autocomplete, AutocompleteRenderInputParams, TextField} from "@mui/material";
import {useQuery} from "@tanstack/react-query";
import {useState} from "react";
import {useDebounce} from "../../../hook/Debounce";

type Props = {
    keys: string[],
    onUpdate: (option: string | null) => void,
    onFetch: (value: string) => Promise<string[]>,
    noOptionsText?: string,
    label?: string,
    placeholder?: string,
    variant?: "standard" | "filled" | "outlined";
    margin?: 'dense' | 'normal' | 'none';
}

export function AutocompleteFetch(props: Props) {
    const {label, placeholder, variant, margin} = props
    const {noOptionsText, keys, onUpdate, onFetch} = props
    const [inputValue, setInputValue] = useState("")
    const debInputValue = useDebounce(inputValue)
    const query = useQuery([...keys, debInputValue], () => onFetch(debInputValue))
    const options = query.data ?? []

    return (
        <Autocomplete
            fullWidth
            renderInput={renderInput}
            inputValue={inputValue}
            noOptionsText={noOptionsText}
            onInputChange={(_, v) => setInputValue(v)}
            options={options}
            loading={query.isLoading}
            onChange={(_, v) => onUpdate(v)}
        />
    )

    function renderInput(params: AutocompleteRenderInputParams) {
        return (
            <TextField
                {...params}
                size={"small"}
                label={label}
                variant={variant}
                margin={margin}
                placeholder={placeholder}
            />
        )
    }
}
