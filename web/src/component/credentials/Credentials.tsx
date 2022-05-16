import {Box, Dialog, DialogContent, DialogTitle, IconButton} from "@mui/material";
import {useStore} from "../../provider/StoreProvider";
import {Close, Security} from "@mui/icons-material";
import React from "react";
import {CredentialsContent} from "./CredentialsContent";

const SX = {
    dialog: { width: "600px" }
}

export function Credentials() {
    const { store, setStore } = useStore()

    return (
        <Dialog open={store.credentialsOpen} onClose={handleClose}>
            <DialogTitle display={"flex"} alignItems={"center"} justifyContent={"space-between"} gap={1}>
                <Security fontSize={"medium"} />
                <Box>Credentials</Box>
                <IconButton size={"small"} onClick={handleClose}><Close /></IconButton>
            </DialogTitle>
            <DialogContent sx={SX.dialog}>
                <CredentialsContent />
            </DialogContent>
        </Dialog>
    )

    function handleClose() {
        setStore({ credentialsOpen: false })
    }
}
