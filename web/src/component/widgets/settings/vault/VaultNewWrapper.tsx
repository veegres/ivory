import {Box} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../../../app/type"
import {AlertCentered} from "../../../view/box/AlertCentered"

const SX: SxPropsMap = {
    box: {padding: "5px 10px 0px 0px", display: "flex", flexDirection: "column", gap: 2},
}

type Props = {
    description: ReactNode,
    children: ReactNode,
}

export function VaultNewWrapper(props: Props) {
    const {description, children} = props
    return (
        <Box sx={SX.box}>
            <AlertCentered text={description}/>
            {children}
        </Box>
    )
}

