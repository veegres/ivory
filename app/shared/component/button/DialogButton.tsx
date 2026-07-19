import {ArrowBack} from "@mui/icons-material"
import {Box, Dialog, DialogActions, DialogTitle, IconButton as MuiIconButton, useMediaQuery, useTheme} from "@mui/material"
import {SvgIconProps} from "@mui/material"
import {ReactElement, ReactNode, useEffect, useState} from "react"

import {SxPropsMap} from "../../helper/HelperType"
import {CloseIconButton} from "./IconButtons"
import {TriggerButton} from "./TriggerButton"

const SX: SxPropsMap = {
    content: {
        width: {xs: "100%", sm: "var(--size-dialog)"}, height: {xs: "auto", sm: "var(--size-dialog)"}, flexGrow: {xs: 1, sm: 0},
        display: "flex", flexDirection: "column",
        gap: 1, padding: "0px 10px 0px 18px ", overflowY: "scroll",
    },
    title: {
        display: "flex", justifyContent: "space-between", alignItems: "center", gap: 1,
        fontFamily: "monospace", padding: "15px 20px 10px"
    },
    action: {display: "flex", justifyContent: "center", gap: 1, padding: "12px 24px"},
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
                    {children}
                </Box>
                <DialogActions sx={SX.action}>
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