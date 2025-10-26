import {Box} from "@mui/material"

import {SxPropsMap} from "../../../app/type"

const SX: SxPropsMap = {
    no: {
        display: "flex", alignItems: "center", justifyContent: "center", textTransform: "uppercase",
        padding: "8px 16px", border: "1px solid", borderColor: "divider", lineHeight: "1.54",
    }
}

type Props = {
    text: string,
}
export function NoBox(props: Props) {
    const {text} = props
    return (
        <Box sx={SX.no}>
            {text}
        </Box>
    )
}
