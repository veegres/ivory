import {Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle} from "@mui/material";

type Props = { open: boolean, title: string, content: string, onAgree: () => void, onClose: () => void }

export function AlertDialog(props: Props) {
    const {open, content, title, onAgree, onClose} = props
    return (
        <Dialog
            open={open}
            onClose={onClose}
            aria-labelledby="alert-dialog-title"
            aria-describedby="alert-dialog-description"
        >
            <DialogTitle id="alert-dialog-title">
                {title}
            </DialogTitle>
            <DialogContent>
                <DialogContentText id="alert-dialog-description">
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
