import {Autocomplete as MuiAutocomplete, AutocompleteRenderInputParams, Box, TextField} from "@mui/material";
import {useMemo, useState} from "react";

export type Option = {
    key: string,
    short: string,
    name: string,
}

type Props = {
    selected: {
        key: string,
        short: string
    },
    options: Option[],
    loading: boolean,
    onUpdate: (option: Option | null) => void,
    label: string,
}

const defaultName = "***"

export function Autocomplete(props: Props) {
    const { onUpdate, loading, label, selected } = props
    const [inputValue, setInputValue] = useState<string>("")
    const {value, options, isOptionNotFound} = useMemo(handleMemoOptions, [selected, props.options])

    return (
        <MuiAutocomplete
            autoHighlight={true}
            options={options}
            value={value}
            onChange={(_, v) => onUpdate(v)}
            inputValue={inputValue}
            onInputChange={(_, v) => setInputValue(v)}
            loading={loading}
            getOptionLabel={(o) => getLabel(o.short, o.name)}
            getOptionDisabled={(o) => o.name === defaultName}
            isOptionEqualToValue={(o, v) => o.key === v.key}
            renderOption={renderOption}
            renderInput={renderInput}
        />
    )

    function renderInput(params: AutocompleteRenderInputParams) {
        return (
            <TextField
                {...params}
                size={"small"}
                label={label}
                error={isOptionNotFound}
                helperText={isOptionNotFound && "element wasn't found"}
            />
        )
    }

    function renderOption(props: React.HTMLAttributes<HTMLLIElement>, option: Option) {
        return (
            <Box component={"li"} {...props}>
                {getLabel(option.short, option.name)}
            </Box>
        )
    }

    function getLabel(shortKey: string, name: string) {
        return `${name} [${shortKey}]`
    }

    function handleMemoOptions() {
        // we have to return `null` here because of Autocomplete has controlled and uncontrolled version
        // depend on `undefined` initial state (https://mui.com/material-ui/react-autocomplete/#controlled-states)
        const selectedOption = props.options.find((option) => option.key === selected.key) ?? null
        const isOptionNotFound = !selectedOption && selected.key !== ""

        // we put options that doesn't exist into the list to avoid problems with Autocomplete component
        if (isOptionNotFound) {
            const notExistedOption = { ...selected, name: defaultName }
            return { value: notExistedOption, options: [notExistedOption, ...props.options], isOptionNotFound }
        } else {
            return { value: selectedOption, options: props.options, isOptionNotFound }
        }
    }
}
