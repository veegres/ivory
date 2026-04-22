import {Box, Paper, ToggleButton, ToggleButtonGroup} from "@mui/material"

import {Node} from "../../../../api/cluster/type"
import {useRouterInfo} from "../../../../api/management/hook"
import {NodeTabType} from "../../../../api/node/type"
import {Permission, PermissionStatus} from "../../../../api/permission/type"
import {ConnectionRequest} from "../../../../api/postgres"
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
    tab: NodeTabType,
    onTab: (tab: NodeTabType) => void,
    connection: ConnectionRequest,
}

export function NodeInfo(props: Props) {
    const {node, tab, onTab, connection} = props
    const info = useRouterInfo(false)
    const permissions = info.data?.auth.user?.permissions
    const access = !!permissions && permissions[Permission.ViewQueryExecuteChart] === PermissionStatus.GRANTED

    return (
        <Box sx={SX.info}>
            <ToggleButtonGroup size={"small"} color={"secondary"} fullWidth value={tab}>
                <ToggleButton value={NodeTabType.CHART} onClick={() => onTab(NodeTabType.CHART)} disabled={!access}>
                    Charts
                </ToggleButton>
                <ToggleButton value={NodeTabType.QUERY} onClick={() => onTab(NodeTabType.QUERY)}>
                    Queries
                </ToggleButton>
            </ToggleButtonGroup>
            <NodeInfoStatus role={node.response.role}/>
            <Paper sx={SX.paper} variant={"outlined"}>
                <NodeInfoTable node={node}/>
            </Paper>
            <Access permission={Permission.ViewQueryExecuteInfo}>
                <Paper sx={SX.paper} variant={"outlined"}>
                    <QueryActivity connection={connection}/>
                </Paper>
            </Access>
        </Box>
    )
}
