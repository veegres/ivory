import {Box, TabScrollButton} from "@mui/material";
import {SxPropsMap} from "../../../app/types";
import {useRef, useState} from "react";
import {ListTagsButton} from "./ListTagsButton";

const ALL = "ALL"
const SCROLL_OFFSET = 100
const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", whiteSpace: "nowrap"},
    arrow: {width: "30px", height: "35px"},
    group: {display: "flex", flexGrow: 1, padding: "0 5px", overflow: "hidden", scrollBehavior: "smooth", gap: 1},
    all: {display: "flex", alignItems: "center", marginRight: "5px"},
    count: {display: "flex", alignItems: "center", marginLeft: "5px"},
}

export function ListTags() {
    const tags = ["test", "test1", "test3"]
    const scrollRef = useRef<Element>(null)
    const [selected, setSelected] = useState(new Set([ALL]))
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
            <Box sx={SX.all}>
                <ListTagsButton tag={ALL} selected={isAll} onClick={handleClickAll}/>
            </Box>
            <Box ref={scrollRef} sx={SX.group}>
                {tags.map(tag => (
                    <ListTagsButton key={tag} tag={tag} selected={selected.has(tag)} onClick={handleClick}/>
                ))}
            </Box>
            <Box sx={SX.count}>
                <ListTagsButton tag={count} selected={!isAll} />
            </Box>
            <TabScrollButton
                sx={SX.arrow}
                direction={"right"}
                orientation={"horizontal"}
                onClick={() => handleScroll(SCROLL_OFFSET)}
            />
        </Box>
    )

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
}
