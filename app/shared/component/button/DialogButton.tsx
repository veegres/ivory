import {ArrowBack} from "@mui/icons-material"
import {Box, Dialog, DialogTitle, IconButton as MuiIconButton, useMediaQuery, useTheme} from "@mui/material"
import {SvgIconProps} from "@mui/material"
import {ReactElement, ReactNode, useState} from "react"

import {SxPropsMap} from "../../helper/HelperType"
import {CloseIconButton} from "./IconButtons"
import {TriggerButton} from "./TriggerButton"

const SX: SxPropsMap = {
    title: {
        display: "flex", justifyContent: "space-between", alignItems: "center", gap: 1,
        // NOTE: no bottom padding - the screen below brings its own top
        // padding, and the two together held the title away from its content
        fontFamily: "monospace", padding: "15px 20px 0px"
    },
}

type Props = {
    title: string
    label?: string,
    children: ReactNode,
    icon: ReactElement<SvgIconProps>,
    size?: number,
    back?: boolean,
    onBackClick?: () => void,
    onClose?: () => void,
    variant?: "button" | "icon" | "button_label",
}

// DialogButton is the shell around a dialog - the trigger, the title bar and
// its navigation. What fills it is one DialogScreen, which brings its own
// content and action bar.
export function DialogButton(props: Props) {
    const {children, title, icon, size, back, onBackClick, onClose, variant = "icon", label} = props
    const [open, setOpen] = useState(false)
    const fullScreen = useMediaQuery(useTheme().breakpoints.down("sm"))

    return (
        <Box>
            {renderTrigger()}
            <Dialog fullScreen={fullScreen} open={open} onClose={handleClose}>
                <DialogTitle sx={SX.title}>
                    <MuiIconButton disableRipple={!back} onClick={onBackClick}>
                        {back ? <ArrowBack/> : icon}
                    </MuiIconButton>
                    <Box>{title}</Box>
                    <CloseIconButton size={40} onClick={handleClose}/>
                </DialogTitle>
                {children}
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

    // NOTE: closing resets the dialog in one step rather than walking its
    // screens back one by one - the arrow is navigation, this is a reset, and
    // driving it from an effect on the closed state re-ran it on every render
    function handleClose() {
        setOpen(false)
        if (onClose) onClose()
    }
}
