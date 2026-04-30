import {Box} from "@mui/material"
import {ReactElement, ReactNode} from "react"

import {SxPropsMap} from "../../../../app/type"

const SX: SxPropsMap = {
    box: {padding: "5px 10px 0px 0px", display: "flex", flexDirection: "column", gap: 2},
    head: {display: "flex", flexDirection: "column", padding: "0px 4px", gap: 0.5},
    title: {display: "flex", alignItems: "center", gap: 1.5, fontSize: 18, fontWeight: 500, color: "text.primary"},
    desc: {fontSize: 12.5, color: "text.secondary", whiteSpace: "pre-wrap", lineHeight: 1.5},
}

type Props = {
    icon: ReactElement,
    label: string,
    description: ReactNode,
    children: ReactNode,
}

export function VaultNewWrapper(props: Props) {
    const {icon, label, description, children} = props
    return (
        <Box sx={SX.box}>
            <Box sx={SX.head}>
                <Box sx={SX.title}>{icon} {label}</Box>
                <Box sx={SX.desc}>{description}</Box>
            </Box>
            {children}
        </Box>
    )
}

