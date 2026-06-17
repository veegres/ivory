import {Box} from "@mui/material"

import {Feature} from "../../../features/feature"
import {ManageAccess} from "../../../features/management/component/ManageAccess"
import {ManageEraseButton} from "../../../features/management/component/ManageEraseButton"
import {ManageFreeButton} from "../../../features/management/component/ManageFreeButton"
import {List} from "../../../shared/component/box/List"
import {ListButton} from "../../../shared/component/box/ListButton"
import {ListItem} from "../../../shared/component/box/ListItem"
import {Settings, SxPropsMap} from "../../../shared/helper/type"
import {SettingOptions} from "../../../shared/helper/utils"
import {ClearCache} from "../browser/ClearCache"
import {SettingsRefetchChanger} from "./SettingsRefetchChanger"
import {SettingsThemeChanger} from "./SettingsThemeChanger"

const SX: SxPropsMap = {
    list: {display: "flex", flexDirection: "column", gap: 3},
}

type Props = {
    onUpdate: (item: Settings) => void
}

export function SettingsContent(props: Props) {
    const {onUpdate} = props

    return (
        <Box sx={SX.list}>
            <List name={"Appearance"}>
                <ListItem title={"Theme"} button={<SettingsThemeChanger/>}/>
                <ListItem title={"Refetch on window focus"} button={<SettingsRefetchChanger/>}/>
            </List>
            <List name={"Privacy and security"}>
                <ManageAccess feature={Feature.ViewVaultList}>{renderButton(Settings.VAULT)}</ManageAccess>
                <ManageAccess feature={Feature.ViewCertList}>{renderButton(Settings.CERTIFICATE)}</ManageAccess>
                <ManageAccess feature={Feature.ManageManagementSecret}>{renderButton(Settings.SECRET)}</ManageAccess>
                {renderButton(Settings.PERMISSION)}
            </List>
            <List name={"Danger Zone"}>
                <ListItem
                    title={"Clear cache"}
                    description={"This clears your local cache. Useful if you experience issues after updates or changes."}
                    button={<ClearCache />}
                />
                <ManageAccess feature={Feature.ManageManagementFree}>
                    <ListItem
                        title={"Free space"}
                        description={"If the VM is running out of disk space, this can help free some space."}
                        button={<ManageFreeButton/>}
                    />
                </ManageAccess>
                <ManageAccess feature={Feature.ManageManagementErase}>
                    <ListItem
                        title={"Erase all data"}
                        description={"Once you erase all data, there is no going back. Please be certain."}
                        button={<ManageEraseButton safe={true}/>}
                    />
                </ManageAccess>
            </List>
            <List name={"About"}>
                <ManageAccess feature={Feature.ManageManagementBackup}>{renderButton(Settings.BACKUP)}</ManageAccess>
                {renderButton(Settings.ABOUT)}
            </List>
        </Box>
    )

    function renderButton(setting: Settings) {
        const {icon, label} = SettingOptions[setting]
        return (
            <ListButton label={label} icon={icon} onClick={() => onUpdate(setting)}/>
        )
    }
}
