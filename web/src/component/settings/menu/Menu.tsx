import {Box, Dialog, DialogContent, DialogTitle} from "@mui/material";
import {useStore} from "../../../provider/StoreProvider";
import {cloneElement, useState} from "react";
import {MenuContent} from "./MenuContent";
import {Settings, SxPropsMap} from "../../../app/types";
import {SettingOptions} from "../../../app/utils";
import {Credentials} from "../credentials/Credentials";
import {Certs} from "../certs/Certs";
import {InfoAlert} from "../../view/InfoAlert";
import {BackIconButton, CloseIconButton} from "../../view/IconButtons";

const SX: SxPropsMap = {
    dialog: {minWidth: "1010px"},
    content: {minWidth: "600px", height: "50vh"},
    title: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 1}
}

export function Menu() {
    const [page, setPage] = useState(Settings.MENU)
    const {store, toggleSettingsDialog} = useStore()
    const options = SettingOptions[page]

    return (
        <Dialog sx={SX.dialog} open={store.settings} onClose={toggleSettingsDialog}>
            <DialogTitle sx={SX.title}>
                {page === Settings.MENU ? cloneElement(options.icon, {sx: {fontSize: 25}}) : (
                    <BackIconButton size={40} onClick={() => setPage(Settings.MENU)}/>
                )}
                <Box>{options.label}</Box>
                <CloseIconButton size={40} onClick={toggleSettingsDialog}/>
            </DialogTitle>
            <DialogContent sx={SX.content}>
                {renderContent()}
            </DialogContent>
        </Dialog>
    )

    function renderContent() {
        switch (page) {
            case Settings.MENU:
                return <MenuContent onUpdate={setPage}/>
            case Settings.PASSWORD:
                return <Credentials/>
            case Settings.CERTIFICATE:
                return <Certs/>
            default:
                return <InfoAlert text={"Not implemented yet"}/>
        }
    }
}
