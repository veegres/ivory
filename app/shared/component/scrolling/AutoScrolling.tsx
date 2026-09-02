import {ArrowDownward, ArrowUpward, Print, RestartAlt, VerticalAlignBottom} from "@mui/icons-material"
import {Box, ToggleButton, Tooltip} from "@mui/material"
import {ReactElement, ReactNode, useEffect, useState} from "react"

import {SxPropsMap} from "../../helper/HelperType"
import {SimpleButton} from "../button/SimpleButton"

const SX: SxPropsMap = {
    wrapper: {display: "flex", flexDirection: "column", gap: 1},
    head: {display: "flex", alignItems: "center", gap: 1},
    title: {
        flexGrow: 1, minWidth: 0, border: 1, borderColor: "divider", padding: "3px 6px", borderRadius: 1,
        fontSize: "12px", fontFamily: "monospace", textTransform: "uppercase", letterSpacing: "0.5px",
        color: "text.secondary", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap",
    },
    buttons: {display: "flex", gap: 0.5},
    button: {minWidth: 0, padding: 0, aspectRatio: "1"},
    icon: {fontSize: "14px"},
}

type Props = {
    auto: boolean,
    length: number,
    scroll: (index: number) => void,
    print?: () => void,
    reconnect?: () => void,
    title?: string,
    children: ReactNode,
}

export function AutoScrolling(props: Props) {
    const {auto, scroll, children, length, print, reconnect, title} = props
    const [autoScrolling, setAutoScrolling] = useState(auto)

    useEffect(handleEffectAutoScrolling, [autoScrolling, length, scroll])
    useEffect(handleEffectSetScrolling, [auto, setAutoScrolling])

    return (
        <Box sx={SX.wrapper}>
            <Box sx={SX.head}>
                <Box sx={SX.title}>{title || "Logs"}</Box>
                <Box sx={SX.buttons}>
                    {renderButtons()}
                </Box>
            </Box>
            {children}
        </Box>
    )

    function renderButtons() {
        return (
            <>
                {reconnect && renderButton("Reconnect", <RestartAlt sx={SX.icon}/>, reconnect)}
                {print && renderButton("Print, Save as PDF", <Print sx={SX.icon}/>, print)}
                {renderToggleButton("Scroll to End", <VerticalAlignBottom sx={SX.icon}/>, () => setAutoScrolling(!autoScrolling))}
                {renderButton("Up the Logs", <ArrowUpward sx={SX.icon}/>, () => {if (length > 0) scroll(0); setAutoScrolling(false)})}
                {renderButton("Down the Logs", <ArrowDownward sx={SX.icon}/>, () => {if (length > 0) scroll(length - 1); setAutoScrolling(false)})}
            </>
        )
    }

    function renderButton(tooltip: string, icon: ReactElement, onClick: () => void) {
        return (
            <Tooltip title={tooltip} placement={"bottom"} arrow={true}>
                <SimpleButton sx={SX.button} color={"inherit"} variant={"outlined"} size={"small"} onClick={onClick}>
                    {icon}
                </SimpleButton>
            </Tooltip>
        )
    }

    function renderToggleButton(tooltip: string, icon: ReactElement, onClick: () => void) {
        return (
            <Tooltip title={tooltip} placement={"bottom"} arrow={true}>
                <ToggleButton sx={SX.button} size={"small"} onClick={onClick} value={"check"} selected={autoScrolling}>
                    {icon}
                </ToggleButton>
            </Tooltip>
        )
    }

    function handleEffectSetScrolling() {
        setAutoScrolling(auto)
    }

    function handleEffectAutoScrolling() {
        if (autoScrolling && length > 0) scroll(length - 1)
    }
}
