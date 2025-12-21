import {Box, Chip} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../../../app/type"

const SX: SxPropsMap = {
    chip: {width: "100%"},
    clusterName: {display: "flex", justifyContent: "center", alignItems: "center", gap: "3px"}
}

type Props = {
    name: string,
    active: boolean,
    onClick: () => void,
    renderRefresh: ReactNode,
}

export function ListCellChip(props: Props) {
    const {name, active, onClick, renderRefresh} = props

    return (
        <Box sx={SX.clusterName}>
            <Chip
                sx={SX.chip}
                color={active ? "primary" : "default"}
                label={name}
                onClick={onClick}
            />
            {renderRefresh}
        </Box>
    )


}
