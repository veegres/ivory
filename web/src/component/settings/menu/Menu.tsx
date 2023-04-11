import {Box, Dialog, DialogContent, DialogTitle, IconButton} from "@mui/material";
import {useStore} from "../../../provider/StoreProvider";
import {useEffect, useState} from "react";
import {MenuContent} from "./MenuContent";
import {Settings, SxPropsMap} from "../../../type/common";
import {SettingOptions} from "../../../app/utils";
import {Credentials} from "../credentials/Credentials";
import {Certs} from "../certs/Certs";
import {InfoAlert} from "../../view/box/InfoAlert";
import {BackIconButton, CloseIconButton} from "../../view/button/IconButtons";
import {About} from "../about/About";
import {Secret} from "../secret/Secret";

const SX: SxPropsMap = {
    dialog: {minWidth: "1010px"},
    content: {minWidth: "600px", height: "60vh"},
    title: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 1},
    menuIcon: { padding: "8px" },
}

export function Menu() {
    const [page, setPage] = useState(Settings.MENU)
    const {store, toggleSettingsDialog} = useStore()
    const options = SettingOptions[page]

    useEffect(handleEffectSettingsClose, [store.settings])

    return (
        <Dialog sx={SX.dialog} open={store.settings} onClose={toggleSettingsDialog}>
            <DialogTitle sx={SX.title}>
                {page === Settings.MENU ? (
                    <IconButton disableRipple>{options.icon}</IconButton>
                ): (
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
            case Settings.SECRET:
                return <Secret/>
            case Settings.ABOUT:
                return <About/>
            default:
                return <InfoAlert text={"Not implemented yet"}/>
        }
    }

    function handleEffectSettingsClose() {
        if (!store.settings) setPage(Settings.MENU)
    }
}
