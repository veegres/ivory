import {ArrowBack} from "@mui/icons-material"
import {Box, Dialog, DialogActions, DialogTitle, IconButton as MuiIconButton, useMediaQuery, useTheme} from "@mui/material"
import {SvgIconProps} from "@mui/material"
import {createContext, ReactElement, ReactNode, useContext, useEffect, useState} from "react"

import {SxPropsMap} from "../../helper/HelperType"
import {CloseIconButton} from "./IconButtons"
import {TriggerButton} from "./TriggerButton"

const SX: SxPropsMap = {
    content: {
        width: {xs: "100%", sm: "var(--size-dialog)"}, maxWidth: "100%", height: {xs: "auto", sm: "var(--size-dialog)"}, flexGrow: {xs: 1, sm: 0},
        display: "flex", flexDirection: "column",
        // NOTE: the top padding is not decoration - without it this scroll
        // container clips the floating label of a first-child text field,
        // which is drawn above its own border
        gap: 1, padding: "10px 10px 0px 18px", overflowY: "scroll",
    },
    title: {
        display: "flex", justifyContent: "space-between", alignItems: "center", gap: 1,
        fontFamily: "monospace", padding: "15px 20px 10px"
    },
    action: {display: "flex", justifyContent: "center", gap: 1, padding: "12px 24px"},
}

// NOTE: undefined means no dialog around us at all, null means the dialog is
// there but its action bar has not attached yet - a child has to tell those
// apart to decide between rendering inline and waiting a tick
const DialogFooterContext = createContext<HTMLElement | null | undefined>(undefined)

// useDialogFooter hands a child the dialog's own action bar, so a form nested
// in the content can put its submit button where every dialog keeps one
// instead of trailing at the end of the scroll.
export function useDialogFooter() {
    return useContext(DialogFooterContext)
}

type Props = {
    title: string
    label?: string,
    children: ReactNode,
    renderActions?: ReactNode,
    icon: ReactElement<SvgIconProps>,
    size?: number,
    back?: boolean,
    onBackClick?: () => void,
    variant?: "button" | "icon" | "button_label",
}

export function DialogButton(props: Props) {
    const {children, renderActions, title, icon, size, back, onBackClick, variant = "icon", label} = props
    const [open, setOpen] = useState(false)
    // NOTE: setFooter is passed as the ref itself - a stable identity, so React
    // does not detach and re-attach it on every render
    const [footer, setFooter] = useState<HTMLElement | null>(null)
    const fullScreen = useMediaQuery(useTheme().breakpoints.down("sm"))

    useEffect(handleEffectClose, [onBackClick, open])

    return (
        <Box>
            {renderTrigger()}
            <Dialog fullScreen={fullScreen} open={open} onClose={() => setOpen(false)}>
                <DialogTitle sx={SX.title}>
                    <MuiIconButton disableRipple={!back} onClick={onBackClick}>
                        {back ? <ArrowBack/> : icon}
                    </MuiIconButton>
                    <Box>{title}</Box>
                    <CloseIconButton size={40} onClick={() => setOpen(false)}/>
                </DialogTitle>
                <Box sx={SX.content}>
                    <DialogFooterContext value={footer}>{children}</DialogFooterContext>
                </Box>
                <DialogActions sx={SX.action} ref={setFooter}>
                    {renderActions}
                </DialogActions>
            </Dialog>
        </Box>
    )

    function renderTrigger() {
        return (
            <TriggerButton
                variant={variant}
                title={title}
                label={label}
                icon={icon}
                size={size}
                onClick={() => setOpen(true)}
            />
        )
    }

    function handleEffectClose() {
        if (!open && onBackClick) onBackClick()
    }
}