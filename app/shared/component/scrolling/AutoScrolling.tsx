import {ArrowDownward, ArrowUpward, Print, RestartAlt, VerticalAlignBottom} from "@mui/icons-material"
import {Box, ToggleButton, Tooltip} from "@mui/material"
import {ReactElement, ReactNode, useEffect, useState} from "react"

import {SxPropsMap} from "../../helper/HelperType"
import {SimpleButton} from "../button/SimpleButton"

const SX: SxPropsMap = {
    wrapper: {display: "flex", gap: 0.5},
    buttons: {display: "flex", flexDirection: "column", gap: 0.5},
    button: {padding: "4px", minWidth: 0},
    icon: {fontSize: "14px"},
}

type Props = {
    auto: boolean,
    length: number,
    scroll: (index: number) => void,
    print?: () => void,
    reconnect?: () => void,
    children: ReactNode,
}

export function AutoScrolling(props: Props) {
    const {auto, scroll, children, length, print, reconnect} = props
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
        return (
            <>
                {renderButton("Up the Logs", <ArrowUpward sx={SX.icon}/>, () => {if (length > 0) scroll(0); setAutoScrolling(false)})}
                {renderButton("Down the Logs", <ArrowDownward sx={SX.icon}/>, () => {if (length > 0) scroll(length - 1); setAutoScrolling(false)})}
                {reconnect && renderButton("Reconnect", <RestartAlt sx={SX.icon}/>, reconnect)}
                {renderToggleButton("Scroll to End", <VerticalAlignBottom sx={SX.icon}/>, () => setAutoScrolling(!autoScrolling))}
                {print && renderButton("Print, Save as PDF", <Print sx={SX.icon}/>, print)}
            </>
        )
    }

    function renderButton(title: string, icon: ReactElement, onClick: () => void) {
        return (
            <Tooltip title={title} placement={"left"} arrow={true}>
                <SimpleButton sx={SX.button} color={"inherit"} variant={"outlined"} size={"small"} onClick={onClick}>
                    {icon}
                </SimpleButton>
            </Tooltip>
        )
    }

    function renderToggleButton(title: string, icon: ReactElement, onClick: () => void) {
        return (
            <Tooltip title={title} placement={"left"} arrow={true}>
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
