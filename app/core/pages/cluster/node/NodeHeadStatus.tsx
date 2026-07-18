import {Box, Skeleton} from "@mui/material"

import {Role} from "../../../../features/node/api/NodeType"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {NodeColor} from "../../../../shared/helper/HelperUtils"

const SX: SxPropsMap = {
    nodeStatusBlock: {
        display: "flex", alignItems: "center", justifyContent: "center",
        minHeight: "70px", maxHeight: "150px", minWidth: "min(300px, 100%)", borderRadius: "4px",
        color: "white", fontSize: "24px", fontWeight: 900, flex: "1 0 min(300px, 100%)",
        boxShadow: "inset 0 0 15px 10px rgb(52 52 52 / 40%)",
    },
}

type Props = {
    role?: Role,
    loading?: boolean,
}

export function NodeHeadStatus(props: Props) {
    const {role, loading} = props
    if (loading) return <Skeleton variant={"rounded"} sx={SX.nodeStatusBlock}/>
    const backgroundColor = role && NodeColor[role].color
    return (
        <Box sx={[SX.nodeStatusBlock, {backgroundColor}]}>
            {role?.toUpperCase() ?? "unknown"}
        </Box>
    )
}
