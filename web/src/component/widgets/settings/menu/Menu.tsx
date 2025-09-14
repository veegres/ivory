import {Box, Dialog, DialogContent, DialogTitle, IconButton} from "@mui/material";
import {useStore, useStoreAction} from "../../../../provider/StoreProvider";
import {useEffect, useState} from "react";
import {MenuContent} from "./MenuContent";
import {SettingOptions} from "../../../../app/utils";
import {Credentials} from "../credentials/Credentials";
import {Certs} from "../certs/Certs";
import {InfoAlert} from "../../../view/box/InfoAlert";
import {BackIconButton, CloseIconButton} from "../../../view/button/IconButtons";
import {About} from "../about/About";
import {Secret} from "../secret/Secret";
import {Settings, SxPropsMap} from "../../../../app/type";

const SX: SxPropsMap = {
    dialog: {minWidth: "1010px"},
    content: {minWidth: "600px", height: "600px"},
    title: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 1},
    menuIcon: { padding: "8px" },
}

export function Menu() {
    const settings = useStore(s => s.settings)
    const {toggleSettingsDialog} = useStoreAction

    const [page, setPage] = useState(Settings.MENU)
    const options = SettingOptions[page]

    useEffect(handleEffectSettingsClose, [settings])

    return (
        <Dialog sx={SX.dialog} open={settings} onClose={toggleSettingsDialog}>
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
        if (!settings) setPage(Settings.MENU)
    }
}
