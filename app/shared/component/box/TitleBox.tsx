import {KeyboardArrowDown} from "@mui/icons-material"
import {Box, Collapse, Typography} from "@mui/material"
import {memo, PropsWithChildren, ReactNode, useState} from "react"

import {SxPropsMap} from "../../helper/HelperType"
import {Hint} from "./Hint"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column"},
    island: {border: 1, borderColor: "divider", borderRadius: 2},
    islandHead: {padding: "6px 8px"},
    islandBody: {padding: "0px 8px 6px"},
    head: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 1, minHeight: "32px"},
    headDense: {minHeight: "26px"},
    label: {flexGrow: 1, minWidth: 0, display: "flex", alignItems: "center", userSelect: "none", gap: 0.5},
    labelToggle: {cursor: "pointer", "&:hover": {color: "primary.main"}, "&:hover *": {color: "primary.main"}},
    labelStatic: {padding: "0px 4px"},
    actions: {display: "flex", alignItems: "center"},
    icon: {fontSize: "20px", transition: "transform 0.2s", transform: "rotate(-90deg)"},
    iconDense: {fontSize: "16px"},
    iconOpen: {transform: "rotate(0deg)"},
    text: {
        fontSize: "15px", fontWeight: 600, textTransform: "uppercase", color: "text.secondary",
        fontFamily: "monospace", transition: "color 0.2s", lineHeight: 1,
    },
    textDense: {fontSize: "13px"},
}

// NOTE: not typed as SxPropsMap - Hint takes a plain SystemStyleObject, and
// the annotation is what makes the two disagree
const HintSX = {
    // NOTE: the padding here depends on haa padding
    hint: {userSelect: "text", cursor: "text", fontSize: "11px", padding: "0px 4px 4px"},
    islandHint: {padding: "0px 12px 4px"},
}

type Props = {
    label: string,
    // NOTE: actions live outside the label so clicking one does not toggle the
    // section it belongs to
    renderActions?: ReactNode,
    defaultOpen?: boolean,
    // NOTE: open/onOpenChange make the section controlled, for when something
    // outside it needs a say - an expand-all button over a list of them. Left
    // out, it keeps its own state as before.
    open?: boolean,
    onOpenChange?: (open: boolean) => void,
    island?: boolean,
    // NOTE: dense is for a section nested inside another box, where the full
    // size label competes with the content it introduces
    dense?: boolean,
    // NOTE: hint sits under the heading, not beside it - it explains the
    // section, so it belongs to the heading rather than to the content, but a
    // sentence long enough to be worth writing never fits on the title's line
    hint?: string,
    // NOTE: collapsible={false} keeps the frame and the heading but drops the
    // toggle, for a section whose content has to stay in view. It exists so
    // such a section is the same box as a collapsible one rather than a
    // hand-built lookalike that drifts from it.
    collapsible?: boolean,
}

export const TitleBox = memo(function SubContentBox(props: PropsWithChildren<Props>) {
    const {label, hint, children, renderActions, defaultOpen = false, island = false, dense = false} = props
    const {collapsible = true, onOpenChange} = props
    const [uncontrolled, setUncontrolled] = useState(defaultOpen)
    const open = props.open ?? uncontrolled

    return (
        <Box sx={[SX.box, island && SX.island]}>
            <Box sx={[SX.head, island && SX.islandHead, dense && SX.headDense]}>
                <Box sx={[SX.label, collapsible ? SX.labelToggle : SX.labelStatic]} onClick={collapsible ? handleToggle : undefined}>
                    {collapsible && renderIcon()}
                    <Typography sx={[SX.text, dense && SX.textDense]}>{label}</Typography>
                </Box>
                {renderActions && <Box sx={SX.actions}>{renderActions}</Box>}
            </Box>
            {renderContent()}
        </Box>
    )

    function renderIcon() {
        return <KeyboardArrowDown sx={[SX.icon, dense && SX.iconDense, open && SX.iconOpen]}/>
    }

    function getHintSx() {
        return {...HintSX.hint, ...(island ? HintSX.islandHint : {})}
    }

    function renderContent() {
        if (!children) return
        if (!collapsible) return renderContentBody()
        return <Collapse in={open}>{renderContentBody()}</Collapse>
    }

    function renderContentBody() {
        return <>
            {hint && <Hint sx={getHintSx()}>{hint}</Hint>}
            <Box sx={[island && SX.islandBody]}>{children}</Box>
        </>
    }

    function handleToggle() {
        setUncontrolled(!open)
        onOpenChange?.(!open)
    }
})
