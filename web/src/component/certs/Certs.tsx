import {Box, Dialog, DialogContent, DialogTitle, IconButton} from "@mui/material";
import {BadgeOutlined, Close} from "@mui/icons-material";
import React from "react";
import {useStore} from "../../provider/StoreProvider";
import {CertsContent} from "./CertsContent";

const SX = {
    content: { width: "600px" },
}

export function Certs() {
    const { store, toggleCertsWindow } = useStore()

    return (
        <Dialog open={store.certsOpen} onClose={toggleCertsWindow}>
            <DialogTitle display={"flex"} alignItems={"center"} justifyContent={"space-between"} gap={1}>
                <BadgeOutlined fontSize={"medium"}/>
                <Box>Certs</Box>
                <IconButton size={"small"} onClick={toggleCertsWindow}><Close/></IconButton>
            </DialogTitle>
            <DialogContent sx={SX.content}>
                <CertsContent/>
            </DialogContent>
        </Dialog>
    )
}
