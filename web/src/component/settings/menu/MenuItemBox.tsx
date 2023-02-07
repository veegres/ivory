import {Box} from "@mui/material";
import {ReactNode} from "react";
import {useTheme} from "../../../provider/ThemeProvider";

const SX = {
    wrapper: {borderRadius: "8px"},
    container: {display: "flex", flexDirection: "column", gap: 1, margin: "0 10px"},
}

type Props = {
    name: string,
    children: ReactNode,
}

export function MenuItemBox(props: Props) {
    const { info } = useTheme()
    const { name, children } = props
    const border = `1px solid ${info?.palette.divider}`

    return (
        <Box sx={SX.container}>
            <Box>{name}</Box>
            <Box sx={{...SX.wrapper, border}}>{children}</Box>
        </Box>
    )
}
