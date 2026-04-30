import {Box, Paper, ToggleButton, ToggleButtonGroup} from "@mui/material"

import {Node} from "../../../../api/cluster/type"
import {Feature} from "../../../../api/feature"
import {useRouterInfo} from "../../../../api/management/hook"
import {NodeTabType} from "../../../../api/node/type"
import {Status} from "../../../../api/permission/type"
import {Connection as QueryConnection} from "../../../../api/query/type"
import {SxPropsMap} from "../../../../app/type"
import {SxPropsFormatter} from "../../../../app/utils"
import {Access} from "../../../widgets/access/Access"
import {QueryActivity} from "../../../widgets/query/QueryActivity"
import {NodeInfoStatus} from "./NodeInfoStatus"
import {NodeInfoTable} from "./NodeInfoTable"

const SX: SxPropsMap = {
    info: {display: "flex", flexDirection: "column", gap: 1, margin: "5px 0", minWidth: "332px"},
    paper: {padding: "5px", bgcolor: SxPropsFormatter.style.bgImageSelected},
}

type Props = {
    node: Node,
    queryCon?: QueryConnection,
    tab: NodeTabType,
    onTab: (tab: NodeTabType) => void,
}

export function NodeInfo(props: Props) {
    const {queryCon, node, tab, onTab} = props
    const info = useRouterInfo(false)
    const permissions = info.data?.auth.user?.permissions
    const access = !!permissions && permissions[Feature.ViewQueryDbChart] === Status.GRANTED

    return (
        <Box sx={SX.info}>
            <ToggleButtonGroup size={"small"} color={"secondary"} fullWidth value={tab}>
                <ToggleButton value={NodeTabType.MONITOR} onClick={() => onTab(NodeTabType.MONITOR)} disabled={!access}>
                    Monitor
                </ToggleButton>
                <ToggleButton value={NodeTabType.QUERY} onClick={() => onTab(NodeTabType.QUERY)}>
                    Queries
                </ToggleButton>
            </ToggleButtonGroup>
            <NodeInfoStatus role={node.keeper.role}/>
            <Paper sx={SX.paper} variant={"outlined"}>
                <NodeInfoTable node={node}/>
            </Paper>
            {renderDbActivity()}
        </Box>
    )

    function renderDbActivity() {
        if (!queryCon) return
        return (
            <Access feature={Feature.ViewQueryDbInfo}>
                <Paper sx={SX.paper} variant={"outlined"}>
                    <QueryActivity connection={queryCon}/>
                </Paper>
            </Access>
        )
    }
}
