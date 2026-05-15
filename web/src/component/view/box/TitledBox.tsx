import {Box} from "@mui/material"
import {PropsWithChildren, ReactNode} from "react"

import {SxPropsMap} from "../../../app/type"

const SX: SxPropsMap = {
    title: {fontSize: "15px", fontWeight: 600, color: "text.secondary", fontFamily: "monospace"},
    head: {display: "flex", justifyContent: "space-between", alignItems: "start", gap: 1, padding: "0px 10px 10px"},
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
        <Box sx={island ? SX.island : {}}>
            <Box sx={SX.head}>
                <Box sx={SX.title}>{title}</Box>
                <Box>{renderActions}</Box>
            </Box>
            {children}
        </Box>
    )
}