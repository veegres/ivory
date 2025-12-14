import {Box} from "@mui/material"

import {Permission} from "../../../../api/permission/type"
import {Settings, SxPropsMap} from "../../../../app/type"
import {SettingOptions} from "../../../../app/utils"
import {List} from "../../../view/box/List"
import {ListButton} from "../../../view/box/ListButton"
import {ListItem} from "../../../view/box/ListItem"
import {Access} from "../../access/Access"
import {ClearCacheButton} from "../../actions/ClearCacheButton"
import {EraseButton} from "../../actions/EraseButton"
import {MenuRefetchChanger} from "./MenuRefetchChanger"
import {MenuThemeChanger} from "./MenuThemeChanger"
import {MenuWrapper} from "./MenuWrapper"

const SX: SxPropsMap = {
    list: {display: "flex", flexDirection: "column", gap: 3},
}

type Props = {
    onUpdate: (item: Settings) => void
}

export function MenuContent(props: Props) {
    const {onUpdate} = props

    return (
        <MenuWrapper>
            <Box sx={SX.list}>
                <List name={"Appearance"}>
                    <ListItem title={"Theme"} button={<MenuThemeChanger/>}/>
                    <ListItem title={"Refetch on window focus"} button={<MenuRefetchChanger/>}/>
                </List>
                <List name={"Privacy and security"}>
                    <Access permission={Permission.ViewPasswordList}>{renderButton(Settings.PASSWORD)}</Access>
                    <Access permission={Permission.ViewCertList}>{renderButton(Settings.CERTIFICATE)}</Access>
                    <Access permission={Permission.ManageManagementSecret}>{renderButton(Settings.SECRET)}</Access>
                    {renderButton(Settings.PERMISSION)}
                </List>
                <List name={"Danger Zone"}>
                    <ListItem
                        title={"Clear cache"}
                        description={`It will clear your local cache. Sometimes it can be helpful when after 
                        updates or some changes you see that something is wrong (counts, selection, etc).`}
                        button={<ClearCacheButton />}
                    />
                    <Access permission={Permission.ManageManagementErase}>
                        <ListItem
                            title={"Erase all data"}
                            description={"Once you erase all data, there is no going back. Please be certain."}
                            button={<EraseButton safe={true}/>}
                        />
                    </Access>
                </List>
                <List name={"About"}>
                    {renderButton(Settings.ABOUT)}
                </List>
            </Box>
        </MenuWrapper>
    )

    function renderButton(setting: Settings) {
        const {icon, label} = SettingOptions[setting]
        return (
            <ListButton label={label} icon={icon} onClick={() => onUpdate(setting)}/>
        )
    }
}
