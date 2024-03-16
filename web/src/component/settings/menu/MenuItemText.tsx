import {Box, ListItem} from "@mui/material";
import {ReactNode} from "react";
import {SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    description: {fontSize: "10px", color: "text.secondary"},
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

    return (
        <ListItem sx={SX.list}>
            <Box sx={SX.body}>
                <Box>{title}</Box>
                <Box sx={SX.description}>{description}</Box>
            </Box>
            {button}
        </ListItem>
    )
}
