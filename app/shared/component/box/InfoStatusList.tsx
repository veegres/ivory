import {Box} from "@mui/material"
import {PropsWithChildren} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    list: {display: "flex", flexWrap: "wrap", gap: "8px 30px", alignItems: "center", width: "100%", padding: "0px 10px"},
    listColumns: {
        display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(140px, max-content))",
        gap: "8px 20px", alignItems: "start", width: "100%", padding: "0px 10px",
    },
    item: {display: "flex", flexDirection: "column", gap: 0.25, fontSize: "0.75rem", fontWeight: "bold", flex: "1 0 70px"},
    itemColumns: {display: "flex", flexDirection: "column", gap: 0.25, fontSize: "0.75rem", fontWeight: "bold", minWidth: 0},
    title: {
        color: "text.disabled", fontSize: "0.7rem", textTransform: "uppercase",
        letterSpacing: "0.05em", whiteSpace: "nowrap",
    },
}

type ListProps = {
    // Lays items out as a grid, so columns stretch to the width of the
    // longest item they contain instead of each item growing independently.
    columns?: boolean,
}

export function InfoStatusList(props: PropsWithChildren<ListProps>) {
    const {columns, children} = props
    return (
        <Box sx={columns ? SX.listColumns : SX.list}>
            {children}
        </Box>
    )
}

type ItemProps = {
    label: string,
    // Must match the columns prop passed to the enclosing InfoStatusList -
    // see the SX.item vs SX.itemColumns comments above for why the two
    // modes cannot share a style.
    columns?: boolean,
}

export function InfoStatusItem(props: PropsWithChildren<ItemProps>) {
    const {label, columns, children} = props
    return (
        <Box sx={columns ? SX.itemColumns : SX.item}>
            <Box sx={SX.title}>{label}</Box>
            {children}
        </Box>
    )
}
