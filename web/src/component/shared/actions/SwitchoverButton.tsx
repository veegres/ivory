import {AlertButton} from "../../view/button/AlertButton";
import {FormControl, InputLabel, MenuItem, Select} from "@mui/material";
import {useState} from "react";
import {Sidecar} from "../../../type/general";
import {InstanceRequest} from "../../../type/instance";
import {useRouterInstanceSwitchover} from "../../../router/instance";
import {Dayjs} from "dayjs";
import {ScheduleInput} from "../../view/input/ScheduleInput";

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
    const body = {leader: request.sidecar.host, candidate, scheduled_at: schedule}

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
                >
                    <MenuItem value={""}><em>none (will be chosen randomly)</em></MenuItem>
                    {candidates.map(sidecar => (
                        <MenuItem key={sidecar.host} value={sidecar.host}>{sidecar.host}</MenuItem>
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
