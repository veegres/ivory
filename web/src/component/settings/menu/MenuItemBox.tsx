import {Box} from "@mui/material";
import {ReactNode} from "react";
import {useTheme} from "../../../provider/ThemeProvider";
import {SxPropsMap} from "../../../app/types";

const SX: SxPropsMap = {
    wrapper: {borderRadius: "8px"},
    container: {display: "flex", flexDirection: "column", gap: 1},
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
