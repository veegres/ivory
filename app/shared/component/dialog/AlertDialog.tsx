import {Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    content: {width: "500px", maxWidth: "100%", display: "flex", gap: 2, flexDirection: "column"},
    title: {width: "500px", maxWidth: "100%", wordBreak: "break-all"},
}

type Props = {
    open: boolean,
    title: string,
    description: ReactNode | string,
    children?: ReactNode,
    onAgree?: () => void,
    onClose: () => void
}

export function AlertDialog(props: Props) {
    const {open, children, description, title, onAgree, onClose} = props
    return (
        <Dialog open={open} onClose={onClose} disableRestoreFocus={true}>
            <DialogTitle sx={SX.title}>
                {title}
            </DialogTitle>
            <DialogContent sx={SX.content}>
                <DialogContentText>
                    {description}
                </DialogContentText>
                {children}
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>{onAgree ? "No, please, No..." : "Close"}</Button>
                {onAgree && <Button onClick={handleAgree} autoFocus>Yes, move on!</Button>}
            </DialogActions>
        </Dialog>
    )

    function handleAgree() {
        if (onAgree) onAgree()
        onClose()
    }
}
