import {Box, TabScrollButton} from "@mui/material"
import {ReactNode, useState} from "react"

import {SxPropsMap} from "../../helper/HelperType"
import {useWindowScrolled} from "../../hook/WindowObservers"

const SCROLL_OFFSET = 100
const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", whiteSpace: "nowrap", minWidth: "0px"},
    arrow: {borderRadius: "5px", margin: "0 3px"},
    group: {display: "flex", flexGrow: 1, overflow: "auto", scrollbarWidth: "none", scrollBehavior: "smooth", gap: 0.5},
    before: {display: "flex", alignItems: "center", marginRight: "5px", lineHeight: "1.1"},
    after: {display: "flex", alignItems: "center", marginLeft: "5px", lineHeight: "1.1"},
}

type Props = {
    children?: ReactNode
    renderBefore?: ReactNode,
    renderAfter?: ReactNode,
    arrowWidth?: string,
    arrowHeight?: string,
    position?: "start" | "end" | "center",
}

export function HiddenScrolling(props: Props) {
    const {children, renderBefore, renderAfter, arrowWidth = "30px", arrowHeight = "35px", position = "start"} = props
    const [ref, setRef] = useState<Element>()
    const [scrolled] = useWindowScrolled(ref)

    return (
        <Box sx={SX.box}>
            <TabScrollButton
                sx={{...SX.arrow, width: arrowWidth, height: arrowHeight}}
                direction={"left"}
                orientation={"horizontal"}
                disabled={!scrolled}
                onClick={() => handleScroll(-SCROLL_OFFSET)}
            />
            {renderBefore && <Box sx={SX.before}>{renderBefore}</Box>}
            {/* NOTE: "safe" keeps the start of centered content reachable when
            it overflows instead of clipping it behind the left arrow */}
            <Box ref={(ref) => setRef(ref as Element)} sx={{...SX.group, justifyContent: position === "start" ? position : `safe ${position}`}}>
                {children}
            </Box>
            {renderAfter && <Box sx={SX.after}>{renderAfter}</Box>}
            <TabScrollButton
                sx={{...SX.arrow, width: arrowWidth, height: arrowHeight}}
                direction={"right"}
                orientation={"horizontal"}
                disabled={!scrolled}
                onClick={() => handleScroll(SCROLL_OFFSET)}
            />
        </Box>
    )

    function handleScroll(scrollOffset: number) {
        const element = ref
        if (element) element.scroll(element.scrollLeft += scrollOffset, 0)
    }
}