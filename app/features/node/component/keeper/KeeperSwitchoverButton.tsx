import {FormControl, InputLabel, MenuItem, Select} from "@mui/material"
import {Dayjs} from "dayjs"
import {useState} from "react"

import {AlertButton} from "../../../../shared/component/button/AlertButton"
import {ScheduleInput} from "../../../../shared/component/input/ScheduleInput"
import {Feature} from "../../../Feature"
import {ManageAccess} from "../../../management/component/ManageAccess"
import {useRouterNodeSwitchover} from "../../api/NodeHook"
import {KeeperCandidate, KeeperOneRequest, Role} from "../../api/NodeType"

type Props = {
    cluster: string,
    request: KeeperOneRequest,
    candidates: KeeperCandidate[],
    role: Role,
    leaderKey?: string,
    size?: "small" | "medium",
}

export function KeeperSwitchoverButton(props: Props) {
    const {request, candidates, cluster, role, leaderKey, size} = props

    const [candidate, setCandidate] = useState<string>()
    const [schedule, setSchedule] = useState<Dayjs>()
    const switchover = useRouterNodeSwitchover(cluster)
    // NOTE: in patroni we cannot use host for leader and candidate, we need to send patroni.name (key)
    const body = {leader: leaderKey, candidate, scheduled_at: schedule}

    return (
        <ManageAccess feature={Feature.ManageNodeKeeperSwitchover}>
            <AlertButton
                size={size}
                color={"secondary"}
                label={"Switchover"}
                title={`Make a switchover of ${request.host}?`}
                description={`It will change the leader of your cluster that will cause some downtime. If you don't choose
                 candidate, the candidate will be chosen randomly.`}
                disabled={role !== "leader"}
                loading={switchover.isPending}
                onClick={handleClick}
            >
                {renderCandidates()}
                <ManageAccess feature={Feature.ManageNodeKeeperSwitchoverSchedule}>
                    <ScheduleInput onChange={(v) => setSchedule(v ?? undefined)} value={schedule ?? null}/>
                </ManageAccess>
            </AlertButton>
        </ManageAccess>
    )

    function renderCandidates() {
        return (
            <FormControl fullWidth>
                <InputLabel id={"select-switchover"}>Candidate</InputLabel>
                <Select
                    labelId={"select-switchover"}
                    label={"Candidate"}
                    value={candidate}
                    onChange={(e) => setCandidate(e.target.value)}
                    fullWidth={true}
                    variant={"outlined"}
                >
                    <MenuItem value={undefined}><em>none (will be chosen randomly)</em></MenuItem>
                    {candidates.map(({name, label}) => (
                        <MenuItem key={name} value={name}>{label}</MenuItem>
                    ))}
                </Select>
            </FormControl>
        )
    }

    function handleClick() {
        switchover.mutate({...request, body})
        setSchedule(undefined)
        setCandidate(undefined)
    }
}
