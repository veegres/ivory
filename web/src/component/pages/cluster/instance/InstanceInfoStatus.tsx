import {Box, Skeleton} from "@mui/material";
import {InstanceColor} from "../../../../app/utils";
import {SxPropsMap} from "../../../../api/management/type";
import {Role} from "../../../../api/instance/type";

const SX: SxPropsMap = {
    instanceStatusBlock: {
        display: "flex", alignItems: "center", justifyContent: "center",
        height: "120px", minWidth: "250px", borderRadius: "4px",
        color: "white", fontSize: "24px", fontWeight: 900,
    },
}

type Props = {
    role?: Role,
    loading?: boolean,
}

export function InstanceInfoStatus(props: Props) {
    const {role, loading} = props
    if (loading) return <Skeleton variant={"rectangular"} sx={SX.instanceStatusBlock}/>
    const background = role && InstanceColor[role]

    return (
        <Box sx={{...SX.instanceStatusBlock, background}}>
            {role?.toUpperCase() ?? "unknown"}
        </Box>
    )
}
