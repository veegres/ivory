import {AlertButton} from "../../view/button/AlertButton";
import {InstanceRequest} from "../../../type/instance";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {InstanceApi} from "../../../app/api";
import {DateTimeField} from "@mui/x-date-pickers";
import {useState} from "react";
import {Box, FormControlLabel, Switch} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {DateTimeFormatter} from "../../../app/utils";

const SX: SxPropsMap = {
    pending: {margin: "0px"},
    description: {color: "text.disabled", fontSize: "12px"},
}

type Props = {
    request: InstanceRequest,
    cluster: string,
}

export function RestartButton(props: Props) {
    const {request, cluster} = props

    const [schedule, setSchedule] = useState<string>()
    const [pending, setPending] = useState(false)

    const options = useMutationOptions([["instance/overview", cluster]])
    const restart = useMutation({mutationFn: InstanceApi.restart, ...options})

    const body = {schedule, restart_pending: pending || undefined}

    return (
        <AlertButton
            size={"small"}
            label={"Restart"}
            title={`Make a restart of ${request.sidecar.host}?`}
            description={"It will restart postgres, that will cause some downtime. Provide local time for schedule."}
            loading={restart.isPending}
            onClick={handleClick}
        >
            <DateTimeField
                label={"Schedule"}
                size={"small"}
                disablePast={true}
                format={DateTimeFormatter.format}
                value={schedule ?? null}
                onChange={(v) => setSchedule(v ?? undefined)}
            />
            <FormControlLabel
                sx={SX.pending}
                disabled={!schedule}
                labelPlacement={"start"}
                control={<Switch checked={pending} onClick={() => setPending(!pending)}/>}
                label={renderLabel()}
            />
        </AlertButton>
    )

    function renderLabel() {
        return (
            <Box>
                <Box>Restart pending</Box>
                <Box sx={SX.description}>
                    will restart Postgres only when restart is pending, requires schedule time
                </Box>
            </Box>
        )
    }

    function handleClick() {
        restart.mutate({...request, body})
        setSchedule(undefined)
        setPending(false)
    }
}
