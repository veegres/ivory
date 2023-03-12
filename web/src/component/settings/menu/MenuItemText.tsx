import {Box, ListItem} from "@mui/material";
import {useTheme} from "../../../provider/ThemeProvider";
import {ReactNode} from "react";
import {SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    description: {fontSize: "10px"},
    list: {display: "flex", justifyContent: "space-between", gap: 1, padding: "5px 16px"},
    body: {display: "flex", flexDirection: "column", margin: "7px 0px"},
}

type Props = {
    title: string,
    description?: string,
    button?: ReactNode,
}

export function MenuItemText(props: Props) {
    const {title, description, button} = props
    const { info } = useTheme()
    const color = info?.palette.text.secondary

    return (
        <ListItem sx={SX.list}>
            <Box sx={SX.body}>
                <Box>{title}</Box>
                <Box sx={{...SX.description, color}}>{description}</Box>
            </Box>
            {button}
        </ListItem>
    )
}
