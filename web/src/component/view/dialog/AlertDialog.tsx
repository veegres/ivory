import {Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {ReactNode} from "react";

const SX: SxPropsMap = {
    dialog: {minWidth: "1010px"},
    content: {width: "500px", display: "flex", gap: 2, flexDirection: "column"},
}

type Props = {
    open: boolean,
    title: string,
    description: string,
    children?: ReactNode,
    onAgree: () => void,
    onClose: () => void
}

export function AlertDialog(props: Props) {
    const {open, children, description, title, onAgree, onClose} = props
    return (
        <Dialog sx={SX.dialog} open={open} onClose={onClose}>
            <DialogTitle>
                {title}
            </DialogTitle>
            <DialogContent sx={SX.content}>
                <DialogContentText>
                    {description}
                </DialogContentText>
                {children}
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>No, please, No...</Button>
                <Button onClick={handleAgree} autoFocus>Yes, move on!</Button>
            </DialogActions>
        </Dialog>
    )

    function handleAgree() {
        onAgree()
        onClose()
    }
}
