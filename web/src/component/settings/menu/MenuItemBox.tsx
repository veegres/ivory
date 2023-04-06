import {Box} from "@mui/material";
import {ReactNode} from "react";
import {SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    wrapper: {
        borderRadius: "8px", border: "1px solid", borderColor: "divider",
        "li:not(:last-child)": {borderBottom: "1px solid", borderColor: "divider"}
    },
    container: {display: "flex", flexDirection: "column", gap: 1},
}

type Props = {
    name: string,
    children: ReactNode,
}

export function MenuItemBox(props: Props) {
    const { name, children } = props

    return (
        <Box sx={SX.container}>
            <Box>{name}</Box>
            <Box sx={SX.wrapper}>{children}</Box>
        </Box>
    )
}
