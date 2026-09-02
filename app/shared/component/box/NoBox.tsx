import {Box} from "@mui/material"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    no: {
        display: "flex", alignItems: "center", justifyContent: "center", textTransform: "uppercase",
        padding: "6px 16px", border: 1, borderColor: "divider", color: "text.secondary", borderRadius: 1,
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
