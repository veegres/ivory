import {SxPropsMap} from "../../app/types";
import {useTheme} from "../../provider/ThemeProvider";
import {useEffect, useRef, useState} from "react";
import {Box, TabScrollButton, ToggleButton} from "@mui/material";

const ALL = "ALL"
const SCROLL_OFFSET = 100
const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", whiteSpace: "nowrap"},
    arrow: {width: "30px", height: "35px"},
    group: {display: "flex", flexGrow: 1, padding: "0 5px", overflow: "hidden", scrollBehavior: "smooth", gap: 1},
    all: {display: "flex", alignItems: "center", marginRight: "5px"},
    count: {display: "flex", alignItems: "center", marginLeft: "5px"},
    element: {padding: "3px 7px", borderRadius: "3px", lineHeight: "1.1"},
}

type Props = {
    tags: string[],
    selected?: string[],
    onUpdate: (tags: string[]) => void
}

export function ToggleButtonScrollable(props: Props) {
    const { tags, onUpdate } = props
    const { info } = useTheme()
    const scrollRef = useRef<Element>(null)
    const [selected, setSelected] = useState(new Set([ALL, ...(props.selected ?? [])]))

    useEffect(handleEffectUpdate, [selected, onUpdate])

    const isAll = selected.has(ALL)
    const count = isAll ? "0" : selected.size.toString()

    return (
        <Box sx={SX.box}>
            <TabScrollButton
                sx={SX.arrow}
                direction={"left"}
                orientation={"horizontal"}
                onClick={() => handleScroll(-SCROLL_OFFSET)}
            />
            <Box sx={SX.all}>{renderButton(ALL, isAll, handleClickAll)}</Box>
            <Box ref={scrollRef} sx={SX.group}>
                {tags.map(tag => renderButton(tag, selected.has(tag), handleClick))}
            </Box>
            <Box sx={SX.count}>{renderButton(count, !isAll)}</Box>
            <TabScrollButton
                sx={SX.arrow}
                direction={"right"}
                orientation={"horizontal"}
                onClick={() => handleScroll(SCROLL_OFFSET)}
            />
        </Box>
    )

    function renderButton(tag: string, selected: boolean, onClick?: (value: string) => void) {
        return (
            <ToggleButton
                sx={{...SX.element, border: `1px solid ${info?.palette.divider}`}}
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
    }

    function handleClickAll() {
        const tmp = new Set([ALL])
        setSelected(tmp)
    }

    function handleScroll(scrollOffset: number) {
        const element = scrollRef.current
        if (element) element.scroll(element.scrollLeft += scrollOffset, 0)
    }

    function handleEffectUpdate() {
        onUpdate([...selected])
    }
}
