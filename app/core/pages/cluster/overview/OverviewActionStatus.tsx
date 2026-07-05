import {Box} from "@mui/material"
import {UseMutationResult} from "@tanstack/react-query"

import {Feature} from "../../../../features/Feature"
import {ManageAccess} from "../../../../features/management/component/ManageAccess"
import {useRouterNodeActivate, useRouterNodePause} from "../../../../features/node/api/NodeHook"
import {KeeperOneRequest, KeeperStatus} from "../../../../features/node/api/NodeType"
import {InfoBox} from "../../../../shared/component/box/InfoBox"
import {AlertButton} from "../../../../shared/component/button/AlertButton"
import {EnumOptions, SxPropsMap} from "../../../../shared/helper/HelperType"
import {KeeperStatusOptions} from "../../../../shared/helper/HelperUtils"

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", gap: 1},
    tooltip: {fontFamily: "monospace"},
}

type Props = {
    cluster: string,
    status: KeeperStatus,
    request: KeeperOneRequest,
}

export function OverviewActionStatus(props: Props) {
    const {status, cluster, request} = props

    const activate = useRouterNodeActivate(cluster)
    const pause = useRouterNodePause(cluster)

    const options = KeeperStatusOptions[status]
    const action: { [key in KeeperStatus]: UseMutationResult<string, any, KeeperOneRequest, unknown> } = {
        [KeeperStatus.Active]: pause,
        [KeeperStatus.Paused]: activate
    }
    const actionButton: { [key in KeeperStatus]: EnumOptions } = {
        [KeeperStatus.Active]: KeeperStatusOptions[KeeperStatus.Paused],
        [KeeperStatus.Paused]: KeeperStatusOptions[KeeperStatus.Active]
    }

    return (
        <Box sx={SX.box}>
            <ManageAccess feature={Feature.ManageNodeKeeperActivation}>
                <AlertButton
                    size={"small"}
                    color={"inherit"}
                    variant={"outlined"}
                    tooltip={<Box sx={SX.tooltip}>{actionButton[status].label}</Box>}
                    label={options.icon}
                    loading={action[status].isPending}
                    title={`Are you sure that you want to ${actionButton[status].label}?`}
                    description={<>This action either active or pause patroni. More info can be
                        found <a href={"https://patroni.readthedocs.io/en/latest/pause.html"}>here</a>.</>}
                    onClick={() => {action[status].mutate(request)}}
                />
            </ManageAccess>
            <InfoBox tooltip={<Box sx={SX.tooltip}>Keeper Status</Box>}>
                <Box sx={{color: options.color}}>{options.name}</Box>
            </InfoBox>
        </Box>
    )
}
