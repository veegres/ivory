import {Box, SxProps, Theme} from "@mui/material"
import {ReactNode, Ref} from "react"

import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {SxPropsFormatter} from "../../../../shared/helper/HelperUtils"

// NOTE: below md the segments are reordered so the name and the actions share
// the first line while the nodes drop to their own full-width line; on md+ the
// visual order and sizes match the old table columns
const SX: SxPropsMap = {
    row: {
        display: "flex", flexWrap: "wrap", alignItems: "flex-start", rowGap: 1, columnGap: "3px", padding: "5px",
        borderBottom: 1, borderColor: "divider", "&:last-child": {borderBottom: 0},
    },
    name: {flex: {xs: "1 1 auto", md: "0 0 210px"}, minWidth: 0, order: 1},
    nodes: {flex: {xs: "1 1 100%", md: "1 1 min(var(--size-input), 100%)"}, minWidth: 0, order: {xs: 3, md: 2}},
    actions: {flex: "0 0 auto", marginLeft: "auto", order: {xs: 2, md: 3}},
}

type Props = {
    renderName: ReactNode,
    renderNodes: ReactNode,
    renderActions: ReactNode,
    sx?: SxProps<Theme>,
    ref?: Ref<HTMLDivElement>,
}

export function ListRowLayout(props: Props) {
    const {renderName, renderNodes, renderActions, sx, ref} = props
    return (
        <Box sx={SxPropsFormatter.merge(SX.row, sx)} ref={ref}>
            <Box sx={SX.name}>{renderName}</Box>
            <Box sx={SX.nodes}>{renderNodes}</Box>
            <Box sx={SX.actions}>{renderActions}</Box>
        </Box>
    )
}
