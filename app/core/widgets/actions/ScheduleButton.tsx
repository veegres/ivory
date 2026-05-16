import {Button} from "@mui/material"

import {Feature} from "../../../features/feature"
import {ScheduledRestart, ScheduledSwitchover} from "../../../features/keeper/type"
import {useRouterNodeRestartDelete, useRouterNodeSwitchoverDelete} from "../../../features/node/hook"
import {KeeperRequest} from "../../../features/node/type"
import {List} from "../../../shared/component/box/List"
import {ListItem} from "../../../shared/component/box/ListItem"
import {NoBox} from "../../../shared/component/box/NoBox"
import {AlertButton} from "../../../shared/component/button/AlertButton"
import {DateTimeFormatter} from "../../../shared/helper/utils"
import {Access} from "../access/Access"

type Props = {
    request: KeeperRequest,
    cluster: string,
    switchover?: ScheduledSwitchover,
    restart?: ScheduledRestart,
}

export function ScheduleButton(props: Props) {
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
            <Access feature={Feature.ManageNodeDbSwitchover}>
                <Button
                    size={"small"}
                    variant={"outlined"}
                    loading={deleteSwitchover.isPending}
                    onClick={() => deleteSwitchover.mutate(request)}
                >
                    Delete
                </Button>
            </Access>
        )
    }

    function renderDeleteRestartButton() {
        return (
            <Access feature={Feature.ManageNodeDbRestart}>
                <Button
                    size={"small"}
                    variant={"outlined"}
                    loading={deleteRestart.isPending}
                    onClick={() => deleteRestart.mutate(request)}
                >
                    Delete
                </Button>
            </Access>
        )
    }
}
