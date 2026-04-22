import {FormControl, InputLabel, MenuItem, Select} from "@mui/material"
import {Dayjs} from "dayjs"
import {useState} from "react"

import {useRouterNodeSwitchover} from "../../../api/node/hook"
import {NodeConnection,NodeRequest} from "../../../api/node/type"
import {Permission} from "../../../api/permission/type"
import {AlertButton} from "../../view/button/AlertButton"
import {ScheduleInput} from "../../view/input/ScheduleInput"
import {Access} from "../access/Access"

type Props = {
    request: NodeRequest,
    cluster: string,
    candidates: NodeConnection[],
    name?: string,
}

export function SwitchoverButton(props: Props) {
    const {request, candidates, cluster, name} = props

    const [candidate, setCandidates] = useState("")
    const [schedule, setSchedule] = useState<Dayjs>()
    const switchover = useRouterNodeSwitchover(cluster)
    // NOTE: in patroni we cannot use host for leader and candidate, we need to send patroni.name
    const body = {leader: name, candidate, scheduled_at: schedule}

    return (
        <Access permission={Permission.ManageNodeDbSwitchover}>
            <AlertButton
                color={"secondary"}
                label={"Switchover"}
                title={`Make a switchover of ${request.connection.host}?`}
                description={`It will change the leader of your cluster that will cause some downtime. If you don't choose
                 candidate, the candidate will be chosen randomly.`}
                loading={switchover.isPending}
                onClick={handleClick}
            >
                {renderCandidates()}
                <ScheduleInput onChange={(v) => setSchedule(v ?? undefined)} value={schedule ?? null}/>
            </AlertButton>
        </Access>
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
                    {candidates.map(connection => (
                        <MenuItem key={connection.host} value={connection.host}>{connection.host}</MenuItem>
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
