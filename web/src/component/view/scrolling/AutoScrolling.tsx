import {Box, Button} from "@mui/material";
import {ReactElement, ReactNode, useEffect, useState} from "react";
import {ArrowDownward, ArrowUpward, Pause, PlayArrow} from "@mui/icons-material";
import {SxPropsMap} from "../../../type/common";
import {SxPropsFormatter} from "../../../app/utils";

const SX: SxPropsMap = {
    wrapper: {display: "flex"},
    buttons: {position: "relative"},
    button: {position: "absolute", right: "15px", zIndex: 1, minWidth: 0},
    top: {top: "10px"},
    bottom: {bottom: "10px"},
}

type Props = {
    auto: boolean,
    length: number,
    scroll: (index: number) => void,
    children: ReactNode,
}

export function AutoScrolling(props: Props) {
    const { auto, scroll, children, length } = props
    const [autoScrolling, setAutoScrolling] = useState(auto)

    useEffect(handleEffectAutoScrolling, [autoScrolling, length, scroll])
    useEffect(handleEffectSetScrolling, [auto, setAutoScrolling])

    return (
        <Box sx={SX.wrapper}>
            {children}
            <Box sx={SX.buttons}>
                {renderButtons()}
            </Box>
        </Box>
    )

    function renderButtons() {
        if (auto) {
            const icon = autoScrolling ? <Pause/> : <PlayArrow/>
            return renderButton("top", icon, () => setAutoScrolling(!autoScrolling))
        }

        return (
            <>
                {renderButton("top", <ArrowUpward/>, () => scroll(0))}
                {renderButton("bottom", <ArrowDownward/>, () => scroll(length - 1))}
            </>
        )
    }

    function renderButton(alignment: "bottom" | "top", icon: ReactElement, onClick: () => void) {
        return (
            <Button sx={SxPropsFormatter.merge(SX.button, SX[alignment])} size={"small"} onClick={onClick}>
                {icon}
            </Button>
        )
    }

    function handleEffectSetScrolling() {
        setAutoScrolling(auto)
    }

    function handleEffectAutoScrolling() {
        if (autoScrolling) scroll(length - 1)
    }
}
