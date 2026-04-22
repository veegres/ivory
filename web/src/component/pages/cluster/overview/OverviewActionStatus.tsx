import {Box} from "@mui/material"
import {UseMutationResult} from "@tanstack/react-query"

import {useRouterNodeActivate, useRouterNodePause} from "../../../../api/node/hook"
import {KeeperStatus,NodeRequest} from "../../../../api/node/type"
import {Permission} from "../../../../api/permission/type"
import {EnumOptions, SxPropsMap} from "../../../../app/type"
import {KeeperStatusOptions} from "../../../../app/utils"
import {InfoBox, Padding} from "../../../view/box/InfoBox"
import {AlertButton} from "../../../view/button/AlertButton"
import {Access} from "../../../widgets/access/Access"

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", gap: 1},
}

type Props = {
    cluster: string,
    status: KeeperStatus,
    request: NodeRequest,
}

export function OverviewActionStatus(props: Props) {
    const {status, cluster, request} = props

    const activate = useRouterNodeActivate(cluster)
    const pause = useRouterNodePause(cluster)

    const options = KeeperStatusOptions[status]
    const action: { [key in KeeperStatus]: UseMutationResult<string, any, NodeRequest, unknown> } = {
        [KeeperStatus.Active]: pause,
        [KeeperStatus.Paused]: activate
    }
    const actionButton: { [key in KeeperStatus]: EnumOptions } = {
        [KeeperStatus.Active]: KeeperStatusOptions.PAUSED,
        [KeeperStatus.Paused]: KeeperStatusOptions.ACTIVE
    }

    return (
        <Box sx={SX.box}>
            <Access permission={Permission.ManageNodeDbActivation}>
                <InfoBox padding={Padding.No}>
                    <AlertButton
                        size={"small"}
                        color={"inherit"}
                        tooltip={actionButton[status].label}
                        label={options.icon}
                        loading={action[status].isPending}
                        title={`Are you sure that you want to ${actionButton[status].label}`}
                        description={<>This action either active or pause patroni. More info can be
                            found <a href={"https://patroni.readthedocs.io/en/latest/pause.html"}>here</a>.</>}
                        onClick={() => {action[status].mutate(request)}}
                    />
                </InfoBox>
            </Access>
            <InfoBox tooltip={"Keeper Status"}>
                <Box sx={{color: options.color}}>{options.name}</Box>
            </InfoBox>
        </Box>
    )
}
