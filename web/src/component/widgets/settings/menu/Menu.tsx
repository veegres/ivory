import {Settings} from "@mui/icons-material"
import {useState} from "react"

import {Settings as SettingsType} from "../../../../app/type"
import {AlertCentered} from "../../../view/box/AlertCentered"
import {DialogButton} from "../../../view/button/DialogButton"
import {About} from "../about/About"
import {Backup} from "../backup/Backup"
import {Certs} from "../certs/Certs"
import {Permissions} from "../permissions/Permissions"
import {Secret} from "../secret/Secret"
import {Vault} from "../vault/Vault"
import {MenuContent} from "./MenuContent"

export function Menu() {
    const [page, setPage] = useState(SettingsType.MENU)

    return (
        <DialogButton
            title={"SETTINGS"}
            renderActions={""}
            icon={<Settings/>}
            back={page !== SettingsType.MENU}
            onBackClick={() => setPage(SettingsType.MENU)}
            size={40}
        >
            {renderContent()}
        </DialogButton>
    )

    function renderContent() {
        switch (page) {
            case SettingsType.MENU:
                return <MenuContent onUpdate={setPage}/>
            case SettingsType.VAULT:
                return <Vault/>
            case SettingsType.CERTIFICATE:
                return <Certs/>
            case SettingsType.SECRET:
                return <Secret/>
            case SettingsType.PERMISSION:
                return <Permissions/>
            case SettingsType.BACKUP:
                return <Backup/>
            case SettingsType.ABOUT:
                return <About/>
            default:
                return <AlertCentered text={"Not implemented yet"}/>
        }
    }
}
