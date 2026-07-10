import {WarningAmber} from "@mui/icons-material"
import {Box} from "@mui/material"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    warning: {
        display: "flex", justifyContent: "center", alignItems: "center",
        color: "warning.main", fontSize: 12, flexWrap: "wrap", gap: 0.5,
    },
}

type Props = {
    warnings: string[],
}

// WarningList shows a list of advisory warnings (e.g. missing values, ignored
// ports); renders nothing when there are none.
export function WarningList(props: Props) {
    const {warnings} = props
    if (warnings.length === 0) return null
    return (
        <Box sx={SX.warning}>
            <WarningAmber sx={{fontSize: 16}}/>
            {warnings.map(warning => <Box key={warning}>{warning}</Box>)}
        </Box>
    )
}
