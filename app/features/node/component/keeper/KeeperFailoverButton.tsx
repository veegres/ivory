import {AlertButton} from "../../../../shared/component/button/AlertButton"
import {Feature} from "../../../feature"
import {ManageAccess} from "../../../management/component/ManageAccess"
import {useRouterNodeFailover} from "../../api/hook"
import {KeeperOneRequest, Role} from "../../api/type"

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
        <ManageAccess feature={Feature.ManageNodeDbFailover}>
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
        </ManageAccess>
    )
}
