import {Autocomplete as MuiAutocomplete, AutocompleteRenderInputParams, Box, TextField} from "@mui/material"
import {HTMLAttributes, useEffect, useMemo, useState} from "react"

const defaultName = "***"

export type Option = { key: string, short: string, name: string }
type Props = {
    selected: {
        key: string,
        short: string
    },
    options: Option[],
    loading: boolean,
    onUpdate: (option: Option | null) => void,
    label: string,
    // search is a persistent baseline query kept in the input while nothing
    // is selected (e.g. a plugin-locked vault username), so it stays visible
    // after blur instead of being cleared
    search?: string,
    // error marks the field from the outside, for a selection a form requires
    // and has not got; the not-found state marks itself
    error?: boolean,
}

export function AutocompleteUuid(props: Props) {
    const {onUpdate, loading, label, selected, search = "", error = false} = props
    const [inputValue, setInputValue] = useState(search)
    const {value, options, isOptionNotFound} = useMemo(handleMemoOptions, [selected, props.options])

    // NOTE: the search baseline can arrive after mount (options load async);
    // sync it into the input while nothing is selected
    useEffect(handleEffectSearch, [search, selected.key])

    return (
        <MuiAutocomplete
            size={"small"}
            autoHighlight={true}
            options={options}
            noOptionsText={"add an option in the settings"}
            value={value}
            onChange={(_, v) => onUpdate(v)}
            inputValue={inputValue}
            onInputChange={handleInputChange}
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
                label={label}
                error={isOptionNotFound || error}
                helperText={isOptionNotFound && "element wasn't found"}
            />
        )
    }

    function renderOption(props: HTMLAttributes<HTMLLIElement>, option: Option) {
        return (
            <Box component={"li"} {...props} key={option.key}>
                {getLabel(option.short, option.name)}
            </Box>
        )
    }

    function handleInputChange(_: unknown, v: string, reason: string) {
        // NOTE: on selection/blur MUI resets the input; when nothing stays
        // selected fall back to the search baseline so it remains visible
        if (reason === "input") setInputValue(v)
        else setInputValue(v || search)
    }

    function handleEffectSearch() {
        if (search && !selected.key) setInputValue(search)
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
            const notExistedOption = {...selected, name: defaultName}
            return {value: notExistedOption, options: [notExistedOption, ...props.options], isOptionNotFound}
        } else {
            return {value: selectedOption, options: props.options, isOptionNotFound}
        }
    }
}
