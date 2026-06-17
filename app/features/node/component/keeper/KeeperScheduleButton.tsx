import {Button} from "@mui/material"

import {List} from "../../../../shared/component/box/List"
import {ListItem} from "../../../../shared/component/box/ListItem"
import {NoBox} from "../../../../shared/component/box/NoBox"
import {AlertButton} from "../../../../shared/component/button/AlertButton"
import {DateTimeFormatter} from "../../../../shared/helper/utils"
import {Feature} from "../../../feature"
import {ManageAccess} from "../../../management/component/ManageAccess"
import {useRouterNodeRestartDelete, useRouterNodeSwitchoverDelete} from "../../api/hook"
import {KeeperOneRequest, ScheduledRestart, ScheduledSwitchover} from "../../api/type"

type Props = {
    request: KeeperOneRequest,
    cluster: string,
    switchover?: ScheduledSwitchover,
    restart?: ScheduledRestart,
}

export function KeeperScheduleButton(props: Props) {
    const {request, cluster, switchover, restart} = props

    const deleteRestart = useRouterNodeRestartDelete(cluster)
    const deleteSwitchover = useRouterNodeSwitchoverDelete(cluster)

    return (
        <AlertButton
            color={"secondary"}
            size={"small"}
            label={"Schedule"}
            title={"Schedule"}
            description={"Here you can check your schedule information and delete it if it is not actual any more."}
            disabled={!switchover && !restart}
            loading={deleteRestart.isPending || deleteSwitchover.isPending}
        >
            <List>
                {restart && (
                    <ListItem
                        title={"Restart"}
                        description={`Scheduled at ${DateTimeFormatter.utc(restart.at)}. Pending restart set to ${restart.pendingRestart}`}
                        button={renderDeleteRestartButton()}
                    />
                )}
                {switchover && (
                    <ListItem
                        title={"Switchover"}
                        description={`Scheduled at ${DateTimeFormatter.utc(switchover.at)}. Candidate set to ${switchover.to}`}
                        button={renderDeleteSwitchoverButton()}
                    />
                )}
            </List>
            {!switchover && !restart && <NoBox text={"There is no schedules yet"}/>}
        </AlertButton>
    )

    function renderDeleteSwitchoverButton() {
        return (
            <ManageAccess feature={Feature.ManageNodeDbSwitchover}>
                <Button
                    size={"small"}
                    variant={"outlined"}
                    loading={deleteSwitchover.isPending}
                    onClick={() => deleteSwitchover.mutate(request)}
                >
                    Delete
                </Button>
            </ManageAccess>
        )
    }

    function renderDeleteRestartButton() {
        return (
            <ManageAccess feature={Feature.ManageNodeDbRestart}>
                <Button
                    size={"small"}
                    variant={"outlined"}
                    loading={deleteRestart.isPending}
                    onClick={() => deleteRestart.mutate(request)}
                >
                    Delete
                </Button>
            </ManageAccess>
        )
    }
}
