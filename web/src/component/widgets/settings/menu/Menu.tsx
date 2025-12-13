import {Box, Dialog, DialogContent, DialogTitle, IconButton} from "@mui/material"
import {useEffect, useState} from "react"

import {Settings, SxPropsMap} from "../../../../app/type"
import {SettingOptions} from "../../../../app/utils"
import {useStore, useStoreAction} from "../../../../provider/StoreProvider"
import {AlertCentered} from "../../../view/box/AlertCentered"
import {BackIconButton, CloseIconButton} from "../../../view/button/IconButtons"
import {About} from "../about/About"
import {Certs} from "../certs/Certs"
import {Credentials} from "../credentials/Credentials"
import {Permissions} from "../permissions/Permissions"
import {Secret} from "../secret/Secret"
import {MenuContent} from "./MenuContent"

const SX: SxPropsMap = {
    dialog: {minWidth: "1010px"},
    content: {minWidth: "600px", height: "600px"},
    title: {display: "flex", alignItems: "center", justifyContent: "space-between", gap: 1},
    menuIcon: {padding: "8px"},
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
            case Settings.PERMISSION:
                return <Permissions/>
            case Settings.ABOUT:
                return <About/>
            default:
                return <AlertCentered text={"Not implemented yet"}/>
        }
    }

    function handleEffectSettingsClose() {
        if (!settings) setPage(Settings.MENU)
    }
}
