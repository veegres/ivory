import {Box} from "@mui/material"

import {SxPropsMap} from "../../../app/type"
import {InfoColorBox} from "./InfoColorBox"

const SX: SxPropsMap = {
    list: {display: "flex", gap: 1, justifyContent: "space-between", margin: "6px 0px 3px 0", maxWidth: "200px", flexWrap: "wrap"},
    label: {textAlign: "center", fontFamily: "monospace"},
}

type Item = { label: string, title?: string, bgColor?: string }
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
                {items.map(({label, title, bgColor}, index) => (
                    <InfoColorBox key={index} label={label} title={title} bgColor={bgColor} opacity={0.9}/>
                ))}
            </Box>
        </Box>
    )
}
