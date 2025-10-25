import {Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle} from "@mui/material";
import {ReactNode} from "react";

import {SxPropsMap} from "../../../app/type";

const SX: SxPropsMap = {
    dialog: {minWidth: "1010px"},
    content: {width: "500px", display: "flex", gap: 2, flexDirection: "column"},
    title: {width: "500px", wordBreak: "break-all"},
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
        <Dialog sx={SX.dialog} open={open} onClose={onClose}>
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
