import {Box} from "@mui/material"
import {UseMutationResult} from "@tanstack/react-query"

import {Feature} from "../../../../features/feature"
import {useRouterNodeActivate, useRouterNodePause} from "../../../../features/node/hook"
import {KeeperOneRequest, KeeperStatus} from "../../../../features/node/type"
import {InfoBox} from "../../../../shared/component/box/InfoBox"
import {AlertButton} from "../../../../shared/component/button/AlertButton"
import {EnumOptions, SxPropsMap} from "../../../../shared/helper/type"
import {KeeperStatusOptions} from "../../../../shared/helper/utils"
import {Access} from "../../../widgets/access/Access"

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
        [KeeperStatus.Active]: KeeperStatusOptions.PAUSED,
        [KeeperStatus.Paused]: KeeperStatusOptions.ACTIVE
    }

    return (
        <Box sx={SX.box}>
            <Access feature={Feature.ManageNodeDbActivation}>
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
            </Access>
            <InfoBox tooltip={<Box sx={SX.tooltip}>Keeper Status</Box>}>
                <Box sx={{color: options.color}}>{options.name}</Box>
            </InfoBox>
        </Box>
    )
}
