import {ListItem, ListItemButton, ListItemIcon, ListItemText} from "@mui/material";
import {NavigateNext} from "@mui/icons-material";
import {ReactElement} from "react";
import {SxPropsMap} from "../../../app/type";

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
