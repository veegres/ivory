import {Box, DialogActions} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    content: {
        width: {xs: "100%", sm: "var(--size-dialog)"}, maxWidth: "100%",
        height: {xs: "auto", sm: "var(--size-dialog)"}, flexGrow: {xs: 1, sm: 0},
        display: "flex", flexDirection: "column",
        gap: 1, padding: "10px 10px 0px 18px", overflowY: "scroll",
    },
    contentAlone: {paddingBottom: "10px"},
    contentFit: {overflowY: "auto", padding: "10px 18px 0px",},
    action: {display: "flex", justifyContent: "center", gap: 1, padding: "12px 24px"},
}

type Props = {
    children: ReactNode,
    renderActions?: ReactNode,
    fit?: boolean,
}

export function DialogScreen(props: Props) {
    const {children, renderActions, fit = false} = props

    return (
        <>
            <Box sx={[SX.content, !renderActions && SX.contentAlone, fit && SX.contentFit]}>{children}</Box>
            {renderActions && <DialogActions sx={SX.action}>{renderActions}</DialogActions>}
        </>
    )
}
