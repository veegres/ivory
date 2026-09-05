import {Box} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../helper/HelperType"
import {AlertCentered} from "./AlertCentered"

const SX: SxPropsMap = {
    text: {display: "flex", flexDirection: "column", gap: 1, textAlign: "center", fontSize: "16px"},
    bold: {fontWeight: "bold", fontSize: "14px"},
    desc: {textAlign: "justify", fontSize: "12px"},
}

type Props = {
    title: ReactNode,
    subtitle: ReactNode,
    description: ReactNode,
}

export function AlertInformative(props: Props) {
    return (
        <AlertCentered text={renderText()}/>
    )

    function renderText() {
        return (
            <Box sx={SX.text}>
                <Box>{props.title}</Box>
                <Box sx={SX.bold}>{props.subtitle}</Box>
                <Box sx={SX.desc}>{props.description}</Box>
            </Box>
        )
    }
}
