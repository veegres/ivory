import {Box} from "@mui/material"
import {ReactNode} from "react"

import {AlertCentered} from "../../../../shared/component/box/AlertCentered"
import {SxPropsMap} from "../../../../shared/helper/type"

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

