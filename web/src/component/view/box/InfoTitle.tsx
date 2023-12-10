import {Box} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {InfoColorBox} from "./InfoColorBox";

const SX: SxPropsMap = {
    box: {display: "flex", gap: 1, justifyContent: "space-between", margin: "3px 0"},
}

type Item = { label: string, title?: string, bgColor?: string }
type Props = {
    items: Item[],
}

export function InfoTitle(props: Props) {
    const {items} = props

    return (
        <Box sx={SX.box}>
            {items.map(({label, title, bgColor}, index) => (
                <InfoColorBox key={index} label={label} title={title} bgColor={bgColor}/>
            ))}
        </Box>
    )
}
