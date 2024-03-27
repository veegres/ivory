import {ListItem, ListItemButton, ListItemIcon, ListItemText} from "@mui/material";
import {NavigateNext} from "@mui/icons-material";
import {SxPropsMap} from "../../../type/general";
import {ReactElement} from "react";

const SX: SxPropsMap = {
    button: {borderRadius: "2px", padding: "12px 16px"},
    label: {margin: "0px"}
}

type Props = {
    label: string,
    icon: ReactElement,
    onClick: () => void,
}

export function ListButton(props: Props) {
    const {label, icon, onClick} = props
    return (
        <ListItem disablePadding>
            <ListItemButton sx={SX.button} onClick={onClick}>
                <ListItemIcon>{icon}</ListItemIcon>
                <ListItemText sx={SX.label}>{label}</ListItemText>
                <NavigateNext/>
            </ListItemButton>
        </ListItem>
    )
}
