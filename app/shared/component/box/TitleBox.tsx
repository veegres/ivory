import {KeyboardArrowDown} from "@mui/icons-material"
import {Box, Collapse, Typography} from "@mui/material"
import {memo, PropsWithChildren, ReactNode, useState} from "react"

import {SxPropsMap} from "../../helper/HelperType"
import {Hint} from "./Hint"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 0.5},
    island: {padding: "6px 8px", border: 1, borderColor: "divider", borderRadius: 2},
    head: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 1},
    label: {flexGrow: 1, minWidth: 0, display: "flex", alignItems: "center", userSelect: "none", gap: 0.5},
    labelToggle: {cursor: "pointer", "&:hover": {color: "primary.main"}, "&:hover *": {color: "primary.main"}},
    actions: {display: "flex", alignItems: "center"},
    icon: {fontSize: "20px", transition: "transform 0.2s", transform: "rotate(-90deg)"},
    iconDense: {fontSize: "16px"},
    iconOpen: {transform: "rotate(0deg)"},
    text: {
        fontSize: "15px", fontWeight: 600, textTransform: "uppercase", color: "text.secondary",
        fontFamily: "monospace", transition: "color 0.2s", lineHeight: 1,
    },
    textDense: {fontSize: "13px"},
    gap: {display: "flex", gap: 0.5, flexDirection: "column"},
}

type Props = {
    label: string,
    renderActions?: ReactNode,
    defaultOpen?: boolean,
    open?: boolean,
    onOpenChange?: (open: boolean) => void,
    island?: boolean,
    dense?: boolean,
    hint?: string,
    collapsible?: boolean,
}

// TitleBox is a labelled section. Every setup lays out the same way - the label,
// the hint under it and the content share one left edge, the toggle follows the
// label instead of pushing it sideways, and island only adds the frame and its
// own padding - so two of them side by side read as one component whatever they
// are given.
export const TitleBox = memo(function TitleBox(props: PropsWithChildren<Props>) {
    const {label, hint, children, renderActions, defaultOpen = false, island = false, dense = false} = props
    const {collapsible = true, onOpenChange} = props
    const [uncontrolled, setUncontrolled] = useState(defaultOpen)
    const open = props.open ?? uncontrolled

    return (
        <Box sx={[SX.box, island && SX.island]}>
            <Box sx={[SX.head, dense && SX.headDense]}>
                <Box sx={[SX.label, collapsible && SX.labelToggle]} onClick={collapsible ? handleToggle : undefined}>
                    <Typography sx={[SX.text, dense && SX.textDense]}>{label}</Typography>
                    {collapsible && renderIcon()}
                </Box>
                {renderActions && <Box sx={SX.actions}>{renderActions}</Box>}
            </Box>
            {renderContent()}
        </Box>
    )

    function renderIcon() {
        return <KeyboardArrowDown sx={[SX.icon, dense && SX.iconDense, open && SX.iconOpen]}/>
    }

    function renderContent() {
        if (!hint && !children) return
        if (!collapsible) return renderContentBody()
        return <Collapse in={open} unmountOnExit={true}>{renderContentBody()}</Collapse>
    }

    function renderContentBody() {
        return (
            <Box sx={SX.gap}>
                {hint && <Hint>{hint}</Hint>}
                {children && <Box>{children}</Box>}
            </Box>
        )
    }

    function handleToggle() {
        setUncontrolled(!open)
        onOpenChange?.(!open)
    }
})
