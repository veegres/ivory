import {KeyboardArrowDown} from "@mui/icons-material"
import {Box, Collapse, Typography} from "@mui/material"
import {memo, PropsWithChildren, ReactNode, useState} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column"},
    // NOTE: less padding vertically than horizontally - the heading already
    // carries its own 4px, so a collapsed box otherwise reads as mostly air
    island: {padding: "4px 8px", border: 1, borderColor: "divider", borderRadius: 2},
    head: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 1},
    label: {
        flexGrow: 1, minWidth: 0,
        display: "flex", alignItems: "center", userSelect: "none",
        gap: 0.5, padding: "4px 0",
    },
    labelToggle: {
        cursor: "pointer", "&:hover": {color: "primary.main"}, "&:hover *": {color: "primary.main"},
    },
    actions: {display: "flex", alignItems: "center"},
    icon: {fontSize: "20px", transition: "transform 0.2s", transform: "rotate(-90deg)"},
    iconDense: {fontSize: "16px"},
    iconOpen: {transform: "rotate(0deg)"},
    // NOTE: only the gap that sits between chevron and text, not the chevron's
    // own width - a heading with no arrow should not be indented as if it had
    // one, it just should not sit flush against the border
    labelStatic: {paddingX: 0.5},
    text: {
        fontSize: "15px", fontWeight: 600, textTransform: "uppercase", color: "text.secondary",
        fontFamily: "monospace", transition: "color 0.2s", lineHeight: 1,
    },
    // NOTE: matches Caption, so a section nested inside a box does not outweigh
    // the fields it contains
    textDense: {fontSize: "13px"},
    // NOTE: a size under Note's - it shares a line with the title rather than
    // sitting under one, so it must not compete with it. It opts back into
    // selection, which the label turns off for the sake of the toggle: a
    // sentence is there to be read and copied, not clicked. Its indent matches
    // labelStatic's, so it starts under the title rather than under the border.
    hint: {
        paddingX: 0.5, userSelect: "text", cursor: "text",
        fontSize: "11px", lineHeight: 1.4, color: "text.disabled",
    },
    content: {marginTop: 1},
    // NOTE: a hint already separates the heading from what follows, so the
    // content takes no gap of its own on top of it
    contentHinted: {marginTop: 0},
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

export const SubContentBox = memo(function SubContentBox(props: PropsWithChildren<Props>) {
    const {label, hint, children, renderActions, defaultOpen = false, island = false, dense = false} = props
    const {collapsible = true, onOpenChange} = props
    const [uncontrolled, setUncontrolled] = useState(defaultOpen)
    const open = props.open ?? uncontrolled

    return (
        <Box sx={[SX.box, island && SX.island]}>
            <Box sx={SX.head}>
                <Box
                    sx={[SX.label, collapsible ? SX.labelToggle : SX.labelStatic]}
                    onClick={collapsible ? handleToggle : undefined}
                >
                    {collapsible && renderIcon()}
                    <Typography sx={[SX.text, dense && SX.textDense]}>{label}</Typography>
                </Box>
                {renderActions && <Box sx={SX.actions}>{renderActions}</Box>}
            </Box>
            {hint && <Box sx={SX.hint}>{hint}</Box>}
            {renderContent()}
        </Box>
    )

    function renderIcon() {
        return <KeyboardArrowDown sx={[SX.icon, dense && SX.iconDense, open && SX.iconOpen]}/>
    }

    function renderContent() {
        const content = <Box sx={[SX.content, !!hint && SX.contentHinted]}>{children}</Box>
        if (!collapsible) return content
        return <Collapse in={open}>{content}</Collapse>
    }

    function handleToggle() {
        setUncontrolled(!open)
        onOpenChange?.(!open)
    }
})
