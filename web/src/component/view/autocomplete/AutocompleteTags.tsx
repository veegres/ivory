import {HTMLAttributes, SyntheticEvent, useEffect, useMemo, useState} from "react";
import {Autocomplete, AutocompleteRenderInputParams, Box, TextField} from "@mui/material";
import {SxPropsMap} from "../../../type/general";

const SX: SxPropsMap = {
    option: {width: "100%", display: "flex", justifyContent: "space-between", alignItems: "center", gap: 1},
    text: {overflowWrap: "anywhere"},
    icon: {borderRadius: "4px", padding: "0 5px", border: "1px solid gray", fontSize: "12px"},
}

type Option = { label: string, exist: boolean }
type TagMap = { [key: string]: boolean }
type Props = {
    tags: string[],
    selected: string[],
    loading: boolean,
    onUpdate: (tags: string[]) => void,
}

export function AutocompleteTags(props: Props) {
    const {tags, selected, onUpdate, loading} = props
    const [inputValue, setInputValue] = useState("")
    const [tagMap, setTagMap] = useState(getTagMap(tags))

    useEffect(handleEffectTagsUpdate, [tags])

    const options = useMemo(() => getOptions(getTagMapWithInput(inputValue, tagMap)), [inputValue, tagMap])
    const value = useMemo(() => getOptions(getTagMap(selected)), [selected])

    return (
        <Autocomplete
            multiple
            size={"small"}
            autoHighlight
            loading={loading}
            options={options}
            value={value}
            noOptionsText={"start typing to add a tag"}
            // NOTE: we need to check is option undefined, because after search when you remove tag it returns
            // undefined, probably this is a bag in mui library
            getOptionLabel={(option) => option?.label ?? ""}
            onInputChange={(_, v) => setInputValue(v.toLowerCase())}
            onChange={handleOnChange}
            isOptionEqualToValue={(o, v) => o.label === v.label}
            filterSelectedOptions
            renderOption={renderOption}
            renderInput={renderInput}
        />
    )

    function renderInput(params: AutocompleteRenderInputParams) {
        return (
            <TextField {...params} label={"Tags"} placeholder={"Tags"}/>
        )
    }

    function renderOption(props: HTMLAttributes<HTMLLIElement>, option: Option) {
        return (
            <Box component={"li"} {...props} key={option.label}>
                <Box sx={SX.option}>
                    <Box sx={SX.text}>{option.label}</Box>
                    {!option.exist && <Box sx={SX.icon}>new</Box>}
                </Box>
            </Box>
        )
    }

    function getTagMapWithInput(input: string, map: TagMap) {
        return input && !map[input] ? {...map, [input]: false} : map
    }

    function getOptions(entries: TagMap) {
        return Object.entries(entries).map(([label, exist]) => ({label, exist}))
    }

    function getTagMap(tags: string[]) {
        return tags.reduce(
            (previousValue, tag) => ({...previousValue, [tag]: true}),
            {} as TagMap
        )
    }

    function handleEffectTagsUpdate() {
        setTagMap(getTagMap(tags))
    }

    function handleOnChange(_o: SyntheticEvent<Element, Event>, v: Option[]) {
        const tags = v.map(v => v.label)
        setTagMap(t => ({...t, ...getTagMap(tags)}))
        onUpdate(tags)
    }
}

