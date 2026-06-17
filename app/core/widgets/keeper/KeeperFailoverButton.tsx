import {Feature} from "../../../features/feature"
import {Role} from "../../../features/keeper/type"
import {useRouterNodeFailover} from "../../../features/node/hook"
import {KeeperOneRequest} from "../../../features/node/type"
import {AlertButton} from "../../../shared/component/button/AlertButton"
import {Access} from "../access/Access"

type Props = {
    name?: string,
    request: KeeperOneRequest,
    cluster: string,
    role: Role,
}

export function KeeperFailoverButton(props: Props) {
    const {request, cluster, role, name} = props

    const failover = useRouterNodeFailover(cluster)
    // NOTE: in patroni we cannot use host for leader and candidate, we need to send patroni.name
    const body = {candidate: name}

    return (
        <Access feature={Feature.ManageNodeDbFailover}>
            <AlertButton
                color={"error"}
                size={"small"}
                label={"Failover"}
                title={`Make a failover of ${request.host}?`}
                description={`It will failover to current node of postgres, that will cause some downtime 
                and potential data loss. Usually it is recommended to use switchover, but if you don't have a
                leader you won't be able to do switchover and here failover can be useful.`}
                disabled={role === "leader"}
                loading={failover.isPending}
                onClick={() => failover.mutate({...request, body})}
            />
        </Access>
    )
}
