import {AlertButton} from "../../view/button/AlertButton";
import {DateTimeField} from "@mui/x-date-pickers";
import {FormControl, InputLabel, MenuItem, Select} from "@mui/material";
import {useState} from "react";
import {Sidecar} from "../../../type/common";
import {InstanceRequest} from "../../../type/instance";
import {useMutation} from "@tanstack/react-query";
import {instanceApi} from "../../../app/api";
import {useMutationOptions} from "../../../hook/QueryCustom";

type Props = {
    request: InstanceRequest,
    cluster: string,
    candidates: Sidecar[],
}

export function SwitchoverButton(props: Props) {
    const {request, candidates, cluster} = props
    const [candidate, setCandidates] = useState("")
    const [schedule, setSchedule] = useState<string | null>(null)

    const options = useMutationOptions([["instance/overview", cluster]])
    const switchover = useMutation({mutationFn: instanceApi.switchover, ...options})

    const body = {leader: request.sidecar.host, candidate, scheduled_at: schedule}

    return (
        <AlertButton
            label={"Switchover"}
            title={`Make a switchover of ${request.sidecar.host}?`}
            description={"It will change the leader of your cluster that will cause some downtime."}
            loading={switchover.isPending}
            onClick={() => switchover.mutate({...request, body})}
            color={"secondary"}
        >
            {renderCandidates()}
            <DateTimeField
                label={"Schedule"}
                size={"small"}
                disablePast={true}
                format={"YYYY-MM-DD HH:mm"}
                value={schedule}
                onChange={(e) => setSchedule(e)}
            />
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
}
