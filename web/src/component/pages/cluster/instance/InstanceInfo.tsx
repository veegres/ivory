import {Box, Paper, ToggleButton, ToggleButtonGroup} from "@mui/material"

import {InstanceTabType, InstanceWeb} from "../../../../api/instance/type"
import {useRouterInfo} from "../../../../api/management/hook"
import {Permission, PermissionStatus} from "../../../../api/permission/type"
import {QueryConnection} from "../../../../api/query/type"
import {SxPropsMap} from "../../../../app/type"
import {SxPropsFormatter} from "../../../../app/utils"
import {Access} from "../../../widgets/access/Access"
import {QueryActivity} from "../../../widgets/query/QueryActivity"
import {InstanceInfoStatus} from "./InstanceInfoStatus"
import {InstanceInfoTable} from "./InstanceInfoTable"

const SX: SxPropsMap = {
    info: {display: "flex", flexDirection: "column", gap: 1, margin: "5px 0", minWidth: "332px"},
    paper: {padding: "5px", bgcolor: SxPropsFormatter.style.bgImageSelected},
}

type Props = {
    instance: InstanceWeb,
    tab: InstanceTabType,
    onTab: (tab: InstanceTabType) => void,
    connection: QueryConnection,
}

export function InstanceInfo(props: Props) {
    const {instance, tab, onTab, connection} = props
    const info = useRouterInfo(false)
    const permissions = info.data?.auth.user?.permissions
    const access = !!permissions && permissions[Permission.ViewQueryExecuteChart] === PermissionStatus.GRANTED

    return (
        <Box sx={SX.info}>
            <ToggleButtonGroup size={"small"} color={"secondary"} fullWidth value={tab}>
                <ToggleButton value={InstanceTabType.CHART} onClick={() => onTab(InstanceTabType.CHART)} disabled={!access}>
                    Charts
                </ToggleButton>
                <ToggleButton value={InstanceTabType.QUERY} onClick={() => onTab(InstanceTabType.QUERY)}>
                    Queries
                </ToggleButton>
            </ToggleButtonGroup>
            <InstanceInfoStatus role={instance.role}/>
            <Paper sx={SX.paper} variant={"outlined"}>
                <InstanceInfoTable instance={instance}/>
            </Paper>
            <Access permission={Permission.ViewQueryExecuteInfo}>
                <Paper sx={SX.paper} variant={"outlined"}>
                    <QueryActivity connection={connection}/>
                </Paper>
            </Access>
        </Box>
    )
}
