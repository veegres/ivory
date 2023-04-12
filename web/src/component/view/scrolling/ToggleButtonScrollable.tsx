import {SxPropsMap} from "../../../type/common";
import {useRef, useState} from "react";
import {Box, TabScrollButton, ToggleButton, Tooltip} from "@mui/material";
import {useWindowScrolled} from "../../../hook/WindowScrolled";

const ALL = "ALL"
const SCROLL_OFFSET = 100
const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", whiteSpace: "nowrap"},
    arrow: {width: "30px", height: "35px", borderRadius: "5px", margin: "0 3px"},
    group: {display: "flex", flexGrow: 1, padding: "0 5px", overflow: "hidden", scrollBehavior: "smooth", gap: 1},
    all: {display: "flex", alignItems: "center", marginRight: "5px", lineHeight: "1.1"},
    info: {display: "flex", alignItems: "center", marginLeft: "5px", lineHeight: "1.1"},
    element: {padding: "3px 7px", borderRadius: "3px", lineHeight: "1.1", border: "1px solid", borderColor: "divider"},
}

type Props = {
    tags: string[],
    selected?: string[],
    warnings?: number,
    onUpdate: (tags: string[]) => void
}

export function ToggleButtonScrollable(props: Props) {
    const {tags, warnings, onUpdate} = props
    const scrollRef = useRef<Element>(null)
    const [scrolled, elements] = useWindowScrolled(scrollRef.current, tags)
    const [selected, setSelected] = useState(new Set([ALL, ...(props.selected ?? [])]))

    const isAll = selected.has(ALL)
    const count = isAll ? "0" : selected.size.toString()

    return (
        <Box sx={SX.box}>
            <TabScrollButton
                sx={SX.arrow}
                direction={"left"}
                orientation={"horizontal"}
                disabled={!scrolled}
                onClick={() => handleScroll(-SCROLL_OFFSET)}
            />
            <Box sx={SX.all}>{renderButton(ALL, isAll, handleClickAll)}</Box>
            <Box ref={scrollRef} sx={SX.group}>
                {elements.map(tag => renderButton(tag, selected.has(tag), handleClick))}
            </Box>
            <Box sx={SX.info}>
                <Tooltip title={"Tags Selected"} placement={"top"}>
                    <span>{renderButton(count, !isAll)}</span>
                </Tooltip>
            </Box>
            <Box sx={SX.info}>
                {renderInfo()}
            </Box>
            <TabScrollButton
                sx={SX.arrow}
                direction={"right"}
                orientation={"horizontal"}
                disabled={!scrolled}
                onClick={() => handleScroll(SCROLL_OFFSET)}
            />
        </Box>
    )

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

    function renderButton(tag: string, selected: boolean, onClick?: (value: string) => void) {
        return (
            <ToggleButton
                sx={SX.element}
                color={"secondary"}
                size={"small"}
                key={tag}
                selected={selected}
                disabled={!onClick}
                value={tag}
                onClick={() => onClick!!(tag)}
            >
                {tag}
            </ToggleButton>
        )
    }

    function handleClick(value: string) {
        const tmp = new Set(selected)
        if (tmp.has(value)) {
            tmp.delete(value)
            if (tmp.size === 0) tmp.add(ALL)
            setSelected(tmp)
        } else {
            tmp.delete(ALL)
            tmp.add(value)
            setSelected(tmp)
        }
        onUpdate([...tmp])
    }

    function handleClickAll() {
        const tmp = new Set([ALL])
        setSelected(tmp)
        onUpdate([...tmp])
    }

    function handleScroll(scrollOffset: number) {
        const element = scrollRef.current
        if (element) element.scroll(element.scrollLeft += scrollOffset, 0)
    }
}