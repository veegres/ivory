import {Box} from "@mui/material"

import {ErrorKeeperRequestMissing} from "../../../../shared/component/box/ErrorManual"
import {TitledBox} from "../../../../shared/component/box/TitledBox"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {Feature} from "../../../Feature"
import {ManageAccess} from "../../../management/component/ManageAccess"
import {KeeperOneRequest, Role} from "../../api/NodeType"
import {KeeperConfig} from "./KeeperConfig"
import {KeeperFailoverButton} from "./KeeperFailoverButton"
import {KeeperReinitButton} from "./KeeperReinitButton"
import {KeeperReloadButton} from "./KeeperReloadButton"
import {KeeperRestartButton} from "./KeeperRestartButton"
import {KeeperScheduleButton} from "./KeeperScheduleButton"
import {KeeperSwitchoverButton} from "./KeeperSwitchoverButton"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", justifyContent: "center", gap: 2},
    actions: {display: "flex", gap: 1, alignItems: "center"},
}

type Props = {
    request?: KeeperOneRequest,
    cluster: string,
    candidates: string[],
    role: Role,
}

export function Keeper(props: Props) {
    const {request, cluster, candidates, role} = props

    return (
        <Box sx={SX.box}>
            <TitledBox title={"Actions"} island={true} renderActions={renderActions()}>
                {!request && <ErrorKeeperRequestMissing/>}
            </TitledBox>
            <ManageAccess feature={Feature.ViewNodeKeeperConfig} error={true}>
                <KeeperConfig req={request}/>
            </ManageAccess>
        </Box>
    )

    function renderActions() {
        if (!request) return
        return (
            <Box sx={SX.actions}>
                <KeeperReloadButton request={request} cluster={cluster}/>
                <KeeperRestartButton request={request} cluster={cluster}/>
                <KeeperReinitButton request={request} cluster={cluster}/>
                <KeeperSwitchoverButton request={request} cluster={cluster} candidates={candidates}/>
                <KeeperFailoverButton request={request} cluster={cluster} role={role}/>
                <KeeperScheduleButton request={request} cluster={cluster}/>
            </Box>
        )
    }
}