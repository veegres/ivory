import {Box} from "@mui/material"

import {SxPropsMap} from "../../helper/type"
import {InfoColorBox} from "./InfoColorBox"

const SX: SxPropsMap = {
    list: {
        display: "flex", gap: 1, margin: "6px 0px 3px 0", padding: "5px", borderRadius: 1,
        maxWidth: "230px", flexWrap: "wrap", bgcolor: "background.default", textTransform: "uppercase",
        opacity: 0.9,
    },
    label: {textAlign: "center", fontFamily: "monospace"},
}

type Item = { label: string, title?: string, color?: string }
type Props = {
    items: Item[],
    label?: string,
}

export function InfoColorBoxList(props: Props) {
    const {items, label} = props

    return (
        <Box>
            {label && <Box sx={SX.label}>{label}</Box>}
            <Box sx={SX.list}>
                {items.map(({label, title, color}, index) => (
                    <InfoColorBox key={index} label={label} title={title} color={color}/>
                ))}
            </Box>
        </Box>
    )
}
