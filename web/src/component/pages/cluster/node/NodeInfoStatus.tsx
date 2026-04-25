import {Box, Skeleton} from "@mui/material"

import {Role} from "../../../../api/keeper/type"
import {SxPropsMap} from "../../../../app/type"
import {NodeColor} from "../../../../app/utils"

const SX: SxPropsMap = {
    nodeStatusBlock: {
        display: "flex", alignItems: "center", justifyContent: "center",
        height: "120px", minWidth: "250px", borderRadius: "4px",
        color: "white", fontSize: "24px", fontWeight: 900,
    },
}

type Props = {
    role?: Role,
    loading?: boolean,
}

export function NodeInfoStatus(props: Props) {
    const {role, loading} = props
    if (loading) return <Skeleton variant={"rectangular"} sx={SX.nodeStatusBlock}/>
    const background = role && NodeColor[role].color

    return (
        <Box sx={{...SX.nodeStatusBlock, background}}>
            {role?.toUpperCase() ?? "unknown"}
        </Box>
    )
}
