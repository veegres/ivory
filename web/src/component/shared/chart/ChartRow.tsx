import {SxPropsMap} from "../../../type/common";
import {Box, Divider} from "@mui/material";
import {ReactNode} from "react";

const SX: SxPropsMap = {
    row: {display: "flex", flexWrap: "wrap", justifyContent: "space-evenly", gap: 3},
}

type Props = {
    label?: ReactNode,
    children: ReactNode,
}

export function ChartRow(props: Props) {
    const {label, children} = props

    return (
        <>
            {label && (<Box><Divider>{label}</Divider></Box>)}
            <Box sx={SX.row}>
                {children}
            </Box>
        </>
    )
}
