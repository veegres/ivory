import {Box, DialogActions} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    content: {
        width: {xs: "100%", sm: "var(--size-dialog)"}, maxWidth: "100%", height: {xs: "auto", sm: "var(--size-dialog)"}, flexGrow: {xs: 1, sm: 0},
        display: "flex", flexDirection: "column",
        // NOTE: the top padding is not decoration - without it this scroll
        // container clips the floating label of a first-child text field,
        // which is drawn above its own border
        gap: 1, padding: "10px 10px 0px 18px", overflowY: "scroll",
    },
    // NOTE: the action bar below is what closes off the content, so a screen
    // without one has to bring its own bottom padding or its last element ends
    // flush against the edge of the dialog
    contentAlone: {paddingBottom: "10px"},
    // NOTE: for a screen that sizes itself to the dialog and scrolls inside its
    // own content - one permanent scrollbar next to another that actually
    // scrolls reads as a broken layout, but the screen still has to be able to
    // scroll where the dialog is short, e.g. full screen on a phone
    contentFit: {overflowY: "auto"},
    action: {display: "flex", justifyContent: "center", gap: 1, padding: "12px 24px"},
}

type Props = {
    children: ReactNode,
    renderActions?: ReactNode,
    fit?: boolean,
}

// DialogScreen is one screen of a DialogButton - its scrolling content plus
// its own action bar. The bar is a sibling of the scroll container rather than
// its last child, which is what keeps it pinned to the bottom of the dialog,
// so a screen renders both halves itself instead of handing its buttons up.
export function DialogScreen(props: Props) {
    const {children, renderActions, fit = false} = props

    return (
        <>
            <Box sx={[SX.content, !renderActions && SX.contentAlone, fit && SX.contentFit]}>{children}</Box>
            {renderActions && <DialogActions sx={SX.action}>{renderActions}</DialogActions>}
        </>
    )
}
