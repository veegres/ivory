import {Alert, Box} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../../app/type"

const SX: SxPropsMap = {
    text: {display: "flex", flexDirection: "column", gap: 1, textAlign: "center"},
    bold: {fontWeight: "bold", fontSize: "14px"},
    desc: {textAlign: "justify", fontSize: "12px", padding: "0px 25px"},
}

type Props = {
    title: ReactNode,
    subtitle: ReactNode,
    description: ReactNode,
}

export function AlertInformative(props: Props) {
    return (
        <Alert severity={"info"} variant={"outlined"} icon={false}>
            <Box sx={SX.text}>
                <Box>{props.title}</Box>
                <Box sx={SX.bold}>{props.subtitle}</Box>
                <Box sx={SX.desc}>{props.description}</Box>
            </Box>
        </Alert>
    )
}
