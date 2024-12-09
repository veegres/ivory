import {Box, CircularProgress, ToggleButton, Tooltip} from "@mui/material";
import {EnumOptions, SidecarStatus, SxPropsMap} from "../../../../type/general";
import {SidecarStatusOptions} from "../../../../app/utils";
import {UseMutationResult} from "@tanstack/react-query";
import {InstanceRequest} from "../../../../type/instance";
import {useRouterInstanceActivate, useRouterInstancePause} from "../../../../router/instance";
import {InfoBox} from "../../../view/box/InfoBox";

const SX: SxPropsMap = {
    box: {display: "flex", alignItems: "center", gap: 1},
    button: {padding: "3px", height: "32px", width: "32px"},
}

type Props = {
    cluster: string,
    status: SidecarStatus,
    request: InstanceRequest,
}

export function OverviewActionStatus(props: Props) {
    const {status, cluster, request} = props

    const activate = useRouterInstanceActivate(cluster)
    const pause = useRouterInstancePause(cluster)

    const options = SidecarStatusOptions[status]
    const action: { [key in SidecarStatus]: UseMutationResult<string, any, InstanceRequest, unknown> } = {
        [SidecarStatus.Active]: pause,
        [SidecarStatus.Paused]: activate
    }
    const actionButton: { [key in SidecarStatus]: EnumOptions } = {
        [SidecarStatus.Active]: SidecarStatusOptions.PAUSED,
        [SidecarStatus.Paused]: SidecarStatusOptions.ACTIVE
    }

    return (
            <Box sx={SX.box}>
                <InfoBox tooltip={"Sidecar Status"} withPadding withRadius>
                    <Box sx={{color: options.color}}>
                        {options.name ?? "UNKNOWN"}
                    </Box>
                </InfoBox>
                <Tooltip title={`${actionButton[status].label} Sidecar`} placement={"top"} arrow={true}>
                    <Box component={"span"}>
                        <ToggleButton
                            sx={SX.button}
                            size={"small"}
                            value={options.label}
                            onClick={() => {action[status].mutate(request)}}
                            disabled={action[status].isPending}
                        >
                            {action[status].isPending ? <CircularProgress size={18} color={"inherit"}/> : options.icon}
                        </ToggleButton>
                    </Box>
                </Tooltip>
            </Box>
    )
}
