import {SxPropsMap} from "../../../type/common";
import {Box} from "@mui/material";
import {ReactNode} from "react";

const SX: SxPropsMap = {
    row: {display: "flex", flexWrap: "wrap", justifyContent: "space-between", gap: 3},
    item: {display: "flex", gap: 1},
    label: {fontWeight: "bold", color: "text.secondary"},
}

type Props = {
    label?: ReactNode,
    children: ReactNode,
}

export function ChartRow(props: Props) {
    const {label, children} = props

    return (
        <>
            {label && (
                <Box sx={SX.item}>
                    <Box sx={SX.label}>Specific stats:</Box>
                    <Box>{label}</Box>
                </Box>
            )}
            <Box sx={SX.row}>
                {children}
            </Box>
        </>
    )
}
