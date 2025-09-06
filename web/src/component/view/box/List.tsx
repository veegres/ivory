import {Box} from "@mui/material";
import {ReactNode} from "react";

import {SxPropsMap} from "../../../app/type";

const SX: SxPropsMap = {
    wrapper: {
        borderRadius: "8px", border: "1px solid", borderColor: "divider",
        "li:not(:last-child)": {borderBottom: "1px solid", borderColor: "divider"}
    },
    container: {display: "flex", flexDirection: "column", gap: 1},
}

type Props = {
    name?: string,
    children: ReactNode,
}

export function List(props: Props) {
    const { name, children } = props

    return (
        <Box sx={SX.container}>
            {name && <Box>{name}</Box>}
            <Box sx={SX.wrapper}>{children}</Box>
        </Box>
    )
}
