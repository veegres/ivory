import {SxPropsMap} from "../../../type/general";
import {useState} from "react";
import {Box, ToggleButton, Tooltip} from "@mui/material";
import {HiddenScrolling} from "./HiddenScrolling";

const ALL = "ALL"
const SX: SxPropsMap = {
    after: {display: "flex", gap: "3px"},
    element: {padding: "3px 7px", borderRadius: "3px", lineHeight: "1.1", border: "1px solid", borderColor: "divider"},
}

type Props = {
    tags: string[],
    selected: string[],
    warnings?: number,
    onUpdate: (tags: string[]) => void
}

export function ToggleButtonScrollable(props: Props) {
    const {tags, selected, warnings, onUpdate} = props
    const tagsSet = new Set(tags)
    const [selectedSet, setSelectedSet] = useState(new Set(selected))

    const isAll = selectedSet.has(ALL)
    const count = isAll ? "0" : selectedSet.size.toString()

    return (
        <HiddenScrolling
            renderBefore={renderSelectButton(ALL, isAll, handleClickAll)}
            renderAfter={renderAfter()}
        >
            {tags.map(tag => renderSelectButton(tag, selectedSet.has(tag), handleClick))}
            {selected.map(tag => renderRemovedButton(tag, !tagsSet.has(tag), handleClick))}
        </HiddenScrolling>
    )

    function renderAfter() {
        return (
            <Box sx={SX.after}>
                <Box sx={SX.info}>
                    <Tooltip title={"Tags Selected"} placement={"top"}>
                        <span>{renderSelectButton(count, !isAll)}</span>
                    </Tooltip>
                </Box>
                <Box sx={SX.info}>
                    {renderInfo()}
                </Box>
            </Box>
        )
    }

    function renderInfo() {
        if (warnings === undefined) return
        return (
            <Tooltip title={"Warnings"} placement={"top"}>
                <span>
                    <ToggleButton
                        sx={SX.element}
                        color={"warning"}
                        size={"small"}
                        selected={warnings > 0}
                        disabled
                        value={warnings}
                    >
                        {warnings}
                    </ToggleButton>
                </span>
            </Tooltip>
        )
    }

    function renderSelectButton(tag: string, selected: boolean, onClick?: (value: string) => void) {
        return (
            <ToggleButton
                sx={SX.element}
                color={"secondary"}
                size={"small"}
                key={tag}
                selected={selected}
                disabled={!onClick}
                value={tag}
                onClick={() => onClick?.(tag)}
            >
                {tag}
            </ToggleButton>
        )
    }

    function renderRemovedButton(tag: string, removed: boolean, onClick?: (value: string) => void) {
        if (tag === ALL || !removed) return null
        return (
            <ToggleButton
                sx={SX.element}
                color={"error"}
                size={"small"}
                key={tag}
                selected={removed}
                disabled={!onClick}
                value={tag}
                onClick={() => onClick?.(tag)}
            >
                {tag}
            </ToggleButton>
        )
    }

    function handleClick(value: string) {
        const tmp = new Set(selectedSet)
        if (tmp.has(value)) {
            tmp.delete(value)
            if (tmp.size === 0) tmp.add(ALL)
            setSelectedSet(tmp)
        } else {
            tmp.delete(ALL)
            tmp.add(value)
            setSelectedSet(tmp)
        }
        onUpdate([...tmp])
    }

    function handleClickAll() {
        const tmp = new Set([ALL])
        setSelectedSet(tmp)
        onUpdate([...tmp])
    }
}
