import {Box, Theme} from "@mui/material"
import {SystemStyleObject} from "@mui/system"
import {PropsWithChildren} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    hint: {color: "text.disabled", fontSize: 12, lineHeight: 1.4, minWidth: 0},
    center: {display: "flex", justifyContent: "center", alignItems: "center", textAlign: "center", gap: 0.5},
}

type Props = {
    center?: boolean,
    sx?: SystemStyleObject<Theme>,
}

// Hint is secondary text that explains or qualifies what it sits next to -
// hints, empty states, descriptions. One size everywhere, so muted text never
// looks like three different things on one screen.
export function Hint(props: PropsWithChildren<Props>) {
    const {center = false, sx, children} = props
    return <Box sx={[SX.hint, center && SX.center, sx ?? {}]}>{children}</Box>
}
