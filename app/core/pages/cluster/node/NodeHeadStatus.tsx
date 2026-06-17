import {Box, Skeleton} from "@mui/material"

import {Role} from "../../../../features/node/api/type"
import {SxPropsMap} from "../../../../shared/helper/type"
import {NodeColor} from "../../../../shared/helper/utils"

const SX: SxPropsMap = {
    nodeStatusBlock: {
        display: "flex", alignItems: "center", justifyContent: "center",
        minHeight: "70px", maxHeight: "150px", minWidth: "300px", borderRadius: "4px",
        color: "white", fontSize: "24px", fontWeight: 900, flex: "1 0 300px",
        boxShadow: "inset 0 0 15px 10px rgb(52 52 52 / 40%)",
    },
}

type Props = {
    role?: Role,
    loading?: boolean,
}

export function NodeHeadStatus(props: Props) {
    const {role, loading} = props
    if (loading) return <Skeleton variant={"rectangular"} sx={SX.nodeStatusBlock}/>
    const backgroundColor = role && NodeColor[role].color
    return (
        <Box sx={[SX.nodeStatusBlock, {backgroundColor}]}>
            {role?.toUpperCase() ?? "unknown"}
        </Box>
    )
}
