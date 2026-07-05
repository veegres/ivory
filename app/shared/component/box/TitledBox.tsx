import {Box} from "@mui/material"
import {PropsWithChildren, ReactNode} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1, backgroundImage: "inherit", backgroundColor: "inherit"},
    title: {fontSize: "1rem", fontWeight: 500, fontFamily: "monospace"},
    head: {display: "flex", justifyContent: "space-between", alignItems: "center", gap: 1, padding: "0px 10px"},
    island: {padding: 1, border: 1, borderColor: "divider", borderRadius: 2},
}

type Props = {
    title: string,
    renderActions?: ReactNode,
    island?: boolean,
}

export function TitledBox(props: PropsWithChildren<Props>) {
    const {title, renderActions, children, island = false} = props
    return (
        <Box sx={island ? [SX.box, SX.island] : SX.box}>
            <Box sx={SX.head}>
                <Box sx={SX.title}>{title}</Box>
                <Box>{renderActions}</Box>
            </Box>
            {children}
        </Box>
    )
}