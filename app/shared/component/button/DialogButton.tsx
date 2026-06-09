import {ArrowBack} from "@mui/icons-material"
import {Box, Dialog, DialogActions, DialogTitle, IconButton as MuiIconButton} from "@mui/material"
import {SvgIconProps} from "@mui/material"
import {ReactElement, ReactNode, useEffect, useState} from "react"

import {SxPropsMap} from "../../helper/type"
import {CloseIconButton, IconButton} from "./IconButtons"

const SX: SxPropsMap = {
    dialog: {minWidth: "1010px"},
    content: {
        width: "600px", height: "600px", display: "flex", flexDirection: "column",
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
    children: ReactNode,
    renderActions: ReactNode,
    icon: ReactElement<SvgIconProps>,
    size?: number,
    back?: boolean,
    onBackClick?: () => void,
}

export function DialogButton(props: Props) {
    const {children, renderActions, title, icon, size, back, onBackClick} = props
    const [open, setOpen] = useState(false)

    useEffect(handleEffectClose, [onBackClick, open])
    
    return (
        <Box>
            <IconButton tooltip={title} icon={icon} size={size} onClick={() => setOpen(true)}/>
            <Dialog sx={SX.dialog} open={open} onClose={() => setOpen(false)}>
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

    function handleEffectClose() {
        if (!open && onBackClick) onBackClick()
    }
}