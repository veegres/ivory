import {Box, TabScrollButton} from "@mui/material";
import {ReactNode, useState} from "react";

import {SxPropsMap} from "../../../app/type";
import {useWindowScrolled} from "../../../hook/WindowObservers";

const SCROLL_OFFSET = 100
const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", whiteSpace: "nowrap"},
    arrow: {borderRadius: "5px", margin: "0 3px"},
    group: {display: "flex", flexGrow: 1, padding: "0 5px", overflow: "hidden", scrollBehavior: "smooth", gap: 1},
    before: {display: "flex", alignItems: "center", marginRight: "5px", lineHeight: "1.1"},
    after: {display: "flex", alignItems: "center", marginLeft: "5px", lineHeight: "1.1"},
}

type Props = {
    children?: ReactNode
    renderBefore?: ReactNode,
    renderAfter?: ReactNode,
    arrowWidth?: string,
    arrowHeight?: string,
}

export function HiddenScrolling(props: Props) {
    const {children, renderBefore, renderAfter, arrowWidth = "30px", arrowHeight = "35px"} = props
    const [ref, setRef] = useState<Element>();
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
            <Box sx={SX.before}>{renderBefore}</Box>
            <Box ref={(ref) => setRef(ref as Element)} sx={SX.group}>{children}</Box>
            <Box sx={SX.after}>{renderAfter}</Box>
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