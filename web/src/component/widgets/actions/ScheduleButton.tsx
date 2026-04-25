import {Button} from "@mui/material"

import {Feature} from "../../../api/feature"
import {useRouterNodeRestartDelete, useRouterNodeSwitchoverDelete} from "../../../api/node/hook"
import {NodeRequest, ScheduledRestart, ScheduledSwitchover} from "../../../api/node/type"
import {DateTimeFormatter} from "../../../app/utils"
import {List} from "../../view/box/List"
import {ListItem} from "../../view/box/ListItem"
import {NoBox} from "../../view/box/NoBox"
import {AlertButton} from "../../view/button/AlertButton"
import {Access} from "../access/Access"

type Props = {
    request: NodeRequest,
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
