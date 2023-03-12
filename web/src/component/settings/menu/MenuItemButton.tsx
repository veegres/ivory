import {SettingOptions} from "../../../app/utils";
import {ListItemButton, ListItemIcon, ListItemText} from "@mui/material";
import {NavigateNext} from "@mui/icons-material";
import {Settings, SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    button: {borderRadius: "8px", padding: "12px 16px"},
    label: {margin: "0px"}
}

type Props = {
    item: Settings,
    onUpdate: (item: Settings) => void,
}

export function MenuItemButton(props: Props) {
    const { item, onUpdate } = props
    const {icon, label} = SettingOptions[item]
    return (
        <ListItemButton sx={SX.button} onClick={() => onUpdate(item)}>
            <ListItemIcon>{icon}</ListItemIcon>
            <ListItemText sx={SX.label}>{label}</ListItemText>
            <NavigateNext/>
        </ListItemButton>
    )
}
