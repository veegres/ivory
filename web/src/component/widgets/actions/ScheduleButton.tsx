import {Button} from "@mui/material"

import {useRouterInstanceRestartDelete, useRouterInstanceSwitchoverDelete} from "../../../api/instance/hook"
import {InstanceRequest, InstanceScheduledRestart, InstanceScheduledSwitchover} from "../../../api/instance/type"
import {Permission} from "../../../api/permission/type"
import {DateTimeFormatter} from "../../../app/utils"
import {List} from "../../view/box/List"
import {ListItem} from "../../view/box/ListItem"
import {NoBox} from "../../view/box/NoBox"
import {AlertButton} from "../../view/button/AlertButton"
import {Access} from "../access/Access"

type Props = {
    request: InstanceRequest,
    cluster: string,
    switchover?: InstanceScheduledSwitchover,
    restart?: InstanceScheduledRestart,
}

export function ScheduleButton(props: Props) {
    const {request, cluster, switchover, restart} = props

    const deleteRestart = useRouterInstanceRestartDelete(cluster)
    const deleteSwitchover = useRouterInstanceSwitchoverDelete(cluster)

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
            <Access permission={Permission.ManageInstanceSwitchover}>
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
            <Access permission={Permission.ManageInstanceRestart}>
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
