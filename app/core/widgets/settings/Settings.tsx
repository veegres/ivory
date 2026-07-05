import {Settings as MuiSettings} from "@mui/icons-material"
import {useState} from "react"

import {Certs} from "../../../features/cert/component/Certs"
import {ManageBackup} from "../../../features/management/component/ManageBackup"
import {ManageSecret} from "../../../features/management/component/ManageSecret"
import {Permissions} from "../../../features/permission/component/Permissions"
import {Vault} from "../../../features/vault/component/Vault"
import {AlertCentered} from "../../../shared/component/box/AlertCentered"
import {DialogButton} from "../../../shared/component/button/DialogButton"
import {Settings as SettingsType} from "../../../shared/helper/HelperType"
import {SettingOptions} from "../../../shared/helper/HelperUtils"
import {SettingsAbout} from "./SettingsAbout"
import {SettingsContent} from "./SettingsContent"

export function Settings() {
    const [page, setPage] = useState(SettingsType.MENU)

    return (
        <DialogButton
            title={SettingOptions[page].label}
            renderActions={""}
            icon={<MuiSettings/>}
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
                return <SettingsContent onUpdate={setPage}/>
            case SettingsType.VAULT:
                return <Vault/>
            case SettingsType.CERTIFICATE:
                return <Certs/>
            case SettingsType.SECRET:
                return <ManageSecret/>
            case SettingsType.PERMISSION:
                return <Permissions/>
            case SettingsType.BACKUP:
                return <ManageBackup/>
            case SettingsType.ABOUT:
                return <SettingsAbout/>
            default:
                return <AlertCentered text={"Not implemented yet"}/>
        }
    }
}
