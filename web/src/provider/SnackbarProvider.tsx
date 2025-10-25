import {Alert, Snackbar} from "@mui/material";
import {createContext, ReactNode, SyntheticEvent, useContext, useEffect, useState} from "react";

interface SnackbarMessage {
    message: string,
    key: number,
    severity: Severity,
}

type Severity = "success" | "info" | "warning" | "error"
type SnackbarType = (message: string, severity: Severity) => void
const SnackbarContext = createContext<SnackbarType>(() => void 0)

export function useSnackbar() {
    return useContext(SnackbarContext)
}

export function SnackbarProvide(props: { children: ReactNode }) {
    const [open, setOpen] = useState(false);
    const [snackPack, setSnackPack] = useState<SnackbarMessage[]>([]);
    const [messageInfo, setMessageInfo] = useState<SnackbarMessage | undefined>();

    useEffect(handleUseEffectStack, [snackPack, messageInfo, open]);

    return (
        <SnackbarContext value={message}>
            {props.children}
            <Snackbar
                key={messageInfo ? messageInfo.key : undefined}
                open={open}
                autoHideDuration={3000}
                onClose={handleClose}
                slotProps={{transition: {onExited: handleExited}}}
            >
                <Alert sx={{width: "100%"}} severity={messageInfo?.severity} variant={"filled"} onClose={handleClose}>
                    {messageInfo ? messageInfo.message : undefined}
                </Alert>
            </Snackbar>
        </SnackbarContext>
    )

    function message(message: string, severity: Severity) {
        setSnackPack((prev) => [...prev, {message, severity, key: new Date().getTime()}])
    }

    function handleClose(_: Event | SyntheticEvent, reason?: string) {
        if (reason === "clickaway") return
        setOpen(false)
    }

    function handleExited() {
        setMessageInfo(undefined)
    }

    function handleUseEffectStack() {
        if (snackPack.length && !messageInfo) {
            // Set a new snack when we don't have an active one
            setMessageInfo({...snackPack[0]});
            setSnackPack((prev) => prev.slice(1));
            setOpen(true);
        } else if (snackPack.length && messageInfo && open) {
            // Close an active snack when a new one is added
            setOpen(false);
        }
    }
}
