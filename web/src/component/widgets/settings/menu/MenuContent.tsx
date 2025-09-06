import {Settings, SxPropsMap} from "../../../../api/management/type";
import {MenuThemeChanger} from "./MenuThemeChanger";
import {EraseButton} from "../../actions/EraseButton";
import {MenuWrapper} from "./MenuWrapper";
import {MenuRefetchChanger} from "./MenuRefetchChanger";
import {Box} from "@mui/material";
import {List} from "../../../view/box/List";
import {ListItem} from "../../../view/box/ListItem";
import {ListButton} from "../../../view/box/ListButton";
import {SettingOptions} from "../../../../app/utils";
import {ClearCacheButton} from "../../actions/ClearCacheButton";
import {MenuUncheckInstanceChanger} from "./MenuUncheckInstanceChanger";

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
                    <ListItem title={"Uncheck instance block on cluster select"} button={<MenuUncheckInstanceChanger/>}/>
                </List>
                <List name={"Privacy and security"}>
                    {renderButton(Settings.PASSWORD)}
                    {renderButton(Settings.CERTIFICATE)}
                    {renderButton(Settings.SECRET)}
                </List>
                <List name={"Danger Zone"}>
                    <ListItem
                        title={"Clear cache"}
                        description={`It will clear your local cache. Sometimes it can be helpful when after 
                        updates or some changes you see that something is wrong (counts, selection, etc).`}
                        button={<ClearCacheButton />}
                    />
                    <ListItem
                        title={"Erase all data"}
                        description={"Once you erase all data, there is no going back. Please be certain."}
                        button={<EraseButton safe={true}/>}
                    />
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
