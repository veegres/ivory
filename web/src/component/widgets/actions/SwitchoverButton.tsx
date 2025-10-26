import {FormControl, InputLabel, MenuItem, Select} from "@mui/material"
import {Dayjs} from "dayjs"
import {useState} from "react"

import {useRouterInstanceSwitchover} from "../../../api/instance/hook"
import {InstanceRequest, Sidecar} from "../../../api/instance/type"
import {AlertButton} from "../../view/button/AlertButton"
import {ScheduleInput} from "../../view/input/ScheduleInput"

type Props = {
    request: InstanceRequest,
    cluster: string,
    candidates: Sidecar[],
}

export function SwitchoverButton(props: Props) {
    const {request, candidates, cluster} = props

    const [candidate, setCandidates] = useState("")
    const [schedule, setSchedule] = useState<Dayjs>()
    const switchover = useRouterInstanceSwitchover(cluster)
    // NOTE: in patroni we cannot use host for leader and candidate, we need to send patroni.name
    const body = {leader: request.sidecar.name, candidate, scheduled_at: schedule}

    return (
        <AlertButton
            color={"secondary"}
            label={"Switchover"}
            title={`Make a switchover of ${request.sidecar.host}?`}
            description={`It will change the leader of your cluster that will cause some downtime. If you don't choose
             candidate, the candidate will be chosen randomly.`}
            loading={switchover.isPending}
            onClick={handleClick}
        >
            {renderCandidates()}
            <ScheduleInput onChange={(v) => setSchedule(v ?? undefined)} value={schedule ?? null}/>
        </AlertButton>
    )

    function renderCandidates() {
        return (
            <FormControl fullWidth size={"small"}>
                <InputLabel id={"select-switchover"}>Candidate</InputLabel>
                <Select
                    labelId={"select-switchover"}
                    label={"Candidate"}
                    value={candidate}
                    onChange={(e) => setCandidates(e.target.value)}
                    fullWidth={true}
                    variant={"outlined"}
                >
                    <MenuItem value={""}><em>none (will be chosen randomly)</em></MenuItem>
                    {candidates.map(sidecar => (
                        <MenuItem key={sidecar.host} value={sidecar.name}>{sidecar.host}</MenuItem>
                    ))}
                </Select>
            </FormControl>
        )
    }

    function handleClick() {
        switchover.mutate({...request, body})
        setSchedule(undefined)
        setCandidates("")
    }
}
