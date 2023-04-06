import {Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle} from "@mui/material";
import {SxPropsMap} from "../../type/common";

const SX: SxPropsMap = {
    dialog: {minWidth: "1010px"},
    content: {width: "500px"},
}

type Props = {
    open: boolean,
    title: string,
    content: string,
    onAgree: () => void,
    onClose: () => void
}

export function AlertDialog(props: Props) {
    const {open, content, title, onAgree, onClose} = props
    return (
        <Dialog sx={SX.dialog} open={open} onClose={onClose}>
            <DialogTitle>
                {title}
            </DialogTitle>
            <DialogContent sx={SX.content}>
                <DialogContentText>
                    {content}
                </DialogContentText>
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
