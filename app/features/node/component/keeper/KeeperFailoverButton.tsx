import {AlertButton} from "../../../../shared/component/button/AlertButton"
import {Feature} from "../../../Feature"
import {ManageAccess} from "../../../management/component/ManageAccess"
import {useRouterNodeFailover} from "../../api/NodeHook"
import {KeeperOneRequest, Role} from "../../api/NodeType"

type Props = {
    name?: string,
    request: KeeperOneRequest,
    cluster: string,
    role: Role,
    size?: "small" | "medium",
}

export function KeeperFailoverButton(props: Props) {
    const {request, cluster, role, name, size} = props

    const failover = useRouterNodeFailover(cluster)
    // NOTE: in patroni we cannot use host for leader and candidate, we need to send patroni.name
    const body = {candidate: name}

    return (
        <ManageAccess feature={Feature.ManageNodeKeeperFailover}>
            <AlertButton
                size={size}
                color={"error"}
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
