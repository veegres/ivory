import {Divider} from "@mui/material";
import {Settings, SxPropsMap} from "../../../app/types";
import {MenuItemBox} from "./MenuItemBox";
import {MenuItemText} from "./MenuItemText";
import {MenuItemButton} from "./MenuItemButton";
import {MenuThemeChanger} from "./MenuThemeChanger";
import {MenuEraseButton} from "./MenuEraseButton";
import {MenuWrapper} from "./MenuWrapper";
import {MenuWrapperScroll} from "./MenuWrapperScroll";
import {MenuRefetchChanger} from "./MenuRefetchChanger";

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
            <MenuWrapperScroll sx={SX.list}>
                <MenuItemBox name={"Appearance"}>
                    <MenuItemText title={"Theme"} button={<MenuThemeChanger/>}/>
                    <Divider/>
                    <MenuItemText title={"Refetch on window focus (coming soon)"} button={<MenuRefetchChanger/>}/>
                </MenuItemBox>
                <MenuItemBox name={"Privacy and security"}>
                    <MenuItemButton item={Settings.PASSWORD} onUpdate={onUpdate}/>
                    <Divider/>
                    <MenuItemButton item={Settings.CERTIFICATE} onUpdate={onUpdate}/>
                    <Divider/>
                    <MenuItemButton item={Settings.SECRET} onUpdate={onUpdate}/>
                </MenuItemBox>
                <MenuItemBox name={"Danger Zone"}>
                    <MenuItemText
                        title={"Erase all data"}
                        description={"Once you erase all data, there is no going back. Please be certain."}
                        button={<MenuEraseButton/>}
                    />
                </MenuItemBox>
                <MenuItemBox name={"About"}>
                    <MenuItemButton item={Settings.ABOUT} onUpdate={onUpdate}/>
                </MenuItemBox>
            </MenuWrapperScroll>
        </MenuWrapper>
    )
}
