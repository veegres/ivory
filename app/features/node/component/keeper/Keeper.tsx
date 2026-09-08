import {Box} from "@mui/material"

import {ErrorKeeperMissing} from "../../../../shared/component/box/ErrorManual"
import {HeadBox} from "../../../../shared/component/box/HeadBox"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {Feature} from "../../../Feature"
import {ManageAccess} from "../../../management/component/ManageAccess"
import {KeeperCandidate, KeeperOneRequest, KeeperOneResponse} from "../../api/NodeType"
import {KeeperConfig} from "./KeeperConfig"
import {KeeperFailoverButton} from "./KeeperFailoverButton"
import {KeeperReinitButton} from "./KeeperReinitButton"
import {KeeperReloadButton} from "./KeeperReloadButton"
import {KeeperRestartButton} from "./KeeperRestartButton"
import {KeeperScheduleButton} from "./KeeperScheduleButton"
import {KeeperSwitchoverButton} from "./KeeperSwitchoverButton"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", justifyContent: "center", gap: 1},
    actions: {display: "flex", gap: 1, alignItems: "center", flexWrap: "wrap"},
}

type Props = {
    request?: KeeperOneRequest,
    cluster: string,
    keeper: KeeperOneResponse,
    candidates: KeeperCandidate[],
}

export function Keeper(props: Props) {
    const {request, cluster, keeper, candidates} = props

    if (!request) return <ErrorKeeperMissing/>

    return (
        <Box sx={SX.box}>
            <HeadBox title={"Actions"}>
                {request && renderActions(request)}
            </HeadBox>
            <ManageAccess feature={Feature.ViewNodeKeeperConfig} error={true}>
                <KeeperConfig req={request}/>
            </ManageAccess>
        </Box>
    )

    function renderActions(request: KeeperOneRequest) {
        return (
            <Box sx={SX.actions}>
                <KeeperReloadButton request={request} cluster={cluster}/>
                <KeeperRestartButton request={request} cluster={cluster}/>
                <KeeperReinitButton request={request} cluster={cluster}/>
                <KeeperSwitchoverButton
                    request={request}
                    cluster={cluster}
                    candidates={candidates}
                    role={keeper.role}
                    leaderKey={keeper.key}
                />
                <KeeperFailoverButton request={request} cluster={cluster} role={keeper.role}/>
                <KeeperScheduleButton
                    request={request}
                    cluster={cluster}
                    switchover={keeper.scheduledSwitchover}
                    restart={keeper.scheduledRestart}
                />
            </Box>
        )
    }
}
