import {Box, Dialog, DialogContent, DialogTitle, IconButton} from "@mui/material";
import {useStore} from "../../../provider/StoreProvider";
import {Close, Security} from "@mui/icons-material";
import React from "react";
import {CredentialsContent} from "./CredentialsContent";

const SX = {
    content: { width: "600px" },
}

export function Credentials() {
    const { store, toggleCredentialsWindow } = useStore()

    return (
        <Dialog open={store.credentialsOpen} onClose={toggleCredentialsWindow}>
            <DialogTitle display={"flex"} alignItems={"center"} justifyContent={"space-between"} gap={1}>
                <Security fontSize={"medium"}/>
                <Box>Credentials</Box>
                <IconButton size={"small"} onClick={toggleCredentialsWindow}><Close/></IconButton>
            </DialogTitle>
            <DialogContent sx={SX.content}>
                <CredentialsContent/>
            </DialogContent>
        </Dialog>
    )
}
