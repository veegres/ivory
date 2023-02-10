import {HTMLAttributes, SyntheticEvent, useEffect, useState} from "react";
import {Autocomplete, AutocompleteRenderInputParams, Box, TextField} from "@mui/material";
import {SxPropsMap} from "../../app/types";

const SX: SxPropsMap = {
    option: {width: "100%", display: "flex", justifyContent: "space-between", alignItems: "center", gap: 1},
    text: {overflowWrap: "anywhere"},
    icon: {borderRadius: "4px", padding: "0 5px", border: "1px solid gray", fontSize: "12px"},
}

type Option = { label: string, exist: boolean }
type TagMap = { [key: string]: boolean }
type Props = {
    tags: string[],
    loading: boolean,
    onUpdate: (tags: string[]) => void,
}

export function AutocompleteTags(props: Props) {
    const {tags, onUpdate, loading} = props
    const [inputValue, setInputValue] = useState("")
    const [objects, setObjects] = useState(getTagsObject(tags))

    useEffect(handleEffectTagsUpdate, [tags])

    const objectsNew = inputValue && !objects[inputValue] ? {...objects, [inputValue]: false} : objects
    const options: Option[] = Object.entries(objectsNew).map(([label, exist]) => ({label, exist}))

    return (
        <Autocomplete
            multiple
            size={"small"}
            autoHighlight
            loading={loading}
            options={options}
            // NOTE: we need to check is option undefined, because after search when you remove tag it returns
            // undefined, probably this is a bag in mui library
            getOptionLabel={(option) => option?.label ?? ""}
            onInputChange={(_, v) => setInputValue(v)}
            onChange={handleOnChange}
            isOptionEqualToValue={(o, v) => o.label === v.label}
            filterSelectedOptions
            renderOption={renderOption}
            renderInput={renderInput}
        />
    )

    function renderInput(params: AutocompleteRenderInputParams) {
        return (
            <TextField{...params} label={"Tags"} placeholder={"Tags"}/>
        )
    }

    function renderOption(props: HTMLAttributes<HTMLLIElement>, option: Option) {
        return (
            <Box component={"li"} {...props}>
                <Box sx={SX.option}>
                    <Box sx={SX.text}>{option.label}</Box>
                    {!option.exist && <Box sx={SX.icon}>new</Box>}
                </Box>
            </Box>
        )
    }

    function getTagsObject(tags: string[]) {
        return tags.reduce(
            (previousValue, tag) => ({...previousValue, [tag]: true}),
            {} as TagMap
        )
    }

    function handleEffectTagsUpdate() {
        setObjects(getTagsObject(tags))
    }

    function handleOnChange(o: SyntheticEvent<Element, Event>, v: Option[]) {
        const tags = v.map(v => v.label)
        setObjects({...objects, ...getTagsObject(tags)})
        onUpdate(tags)
    }
}

