import {alpha, Box, Skeleton, Theme} from "@mui/material"

import {Role} from "../../../../features/node/api/NodeType"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {NodeColor} from "../../../../shared/helper/HelperUtils"

const SX: SxPropsMap = {
    block: {
        display: "flex", alignItems: "center", justifyContent: "center",
        minHeight: "70px", minWidth: "min(var(--size-field), 100%)", alignSelf: "stretch",
        flex: "1 0 min(var(--size-field), 100%)", border: 1, borderRadius: 2,
        fontFamily: "monospace", fontSize: "22px", fontWeight: 600, letterSpacing: "2px",
        textTransform: "uppercase", userSelect: "none",
    },
}

type Props = {
    role?: Role,
    loading?: boolean,
}

export function NodeHeadRole(props: Props) {
    const {role, loading} = props
    if (loading) return <Skeleton variant={"rounded"} sx={SX.block}/>
    const {label} = NodeColor[role ?? "unknown"]
    return (
        <Box sx={[SX.block, getSurface]}>
            {role ?? "unknown"}
        </Box>
    )

    function getSurface(theme: Theme) {
        const {main} = theme.palette[label]
        return {
            color: main,
            borderColor: alpha(main, 0.35),
            background: `linear-gradient(145deg, ${alpha(main, 0.16)} 0%, ${alpha(main, 0.05)} 100%)`,
        }
    }
}
