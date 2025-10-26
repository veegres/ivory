import {Box, FormControlLabel, Switch} from "@mui/material"
import {Dayjs} from "dayjs"
import {useState} from "react"

import {useRouterInstanceRestart} from "../../../api/instance/hook"
import {InstanceRequest} from "../../../api/instance/type"
import {SxPropsMap} from "../../../app/type"
import {AlertButton} from "../../view/button/AlertButton"
import {ScheduleInput} from "../../view/input/ScheduleInput"

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

    const [schedule, setSchedule] = useState<Dayjs>()
    const [pending, setPending] = useState(false)

    const restart = useRouterInstanceRestart(cluster)

    const body = {schedule, restart_pending: pending || undefined}

    return (
        <AlertButton
            size={"small"}
            label={"Restart"}
            title={`Make a restart of ${request.sidecar.host}?`}
            description={"It will restart postgres, that will cause some downtime."}
            loading={restart.isPending}
            onClick={handleClick}
        >
            <ScheduleInput onChange={(v) => setSchedule(v ?? undefined)} value={schedule ?? null}/>
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
                    Will restart Postgres only when restart is pending, requires schedule time is set.
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
