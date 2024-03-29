import {Box, CircularProgress, ToggleButton, Tooltip} from "@mui/material";
import {SidecarStatus, SxPropsMap} from "../../../../type/general";
import {SidecarStatusOptions} from "../../../../app/utils";
import {InfoColorBoxList} from "../../../view/box/InfoColorBoxList";
import {UseMutationResult} from "@tanstack/react-query";
import {InstanceRequest} from "../../../../type/instance";
import {useRouterInstanceActivate, useRouterInstancePause} from "../../../../router/instance";

const SX: SxPropsMap = {
    button: {padding: "3px", height: "32px", width: "32px"},
}

type Props = {
    cluster: string,
    status: SidecarStatus,
    instance: InstanceRequest,
}

export function OverviewActionStatus(props: Props) {
    const {status, cluster, instance} = props

    const activate = useRouterInstanceActivate(cluster)
    const pause = useRouterInstancePause(cluster)

    const options = SidecarStatusOptions[status]
    const item = {label: options.label, bgColor: options.color}
    const action: { [key in SidecarStatus]: UseMutationResult<string, any, InstanceRequest, unknown> }= {
        [SidecarStatus.Active]: pause,
        [SidecarStatus.Paused]: activate
    }

    return (
        <Tooltip title={<InfoColorBoxList items={[item]} label={"Sidecar Status"}/>} placement={"top"} arrow={true}>
            <Box component={"span"}>
                <ToggleButton
                    sx={SX.button}
                    size={"small"}
                    value={options.label}
                    onClick={() => { action[status].mutate(instance) }}
                    disabled={action[status].isPending}
                >
                    {action[status].isPending ? <CircularProgress size={18} color={"inherit"}/> : options.icon}
                </ToggleButton>
            </Box>
        </Tooltip>
    )
}
