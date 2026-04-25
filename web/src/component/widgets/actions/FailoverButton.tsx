import {Feature} from "../../../api/feature"
import {useRouterNodeFailover} from "../../../api/node/hook"
import {KeeperRequest} from "../../../api/node/type"
import {AlertButton} from "../../view/button/AlertButton"
import {Access} from "../access/Access"

type Props = {
    request: KeeperRequest,
    cluster: string,
    disabled: boolean,
    name?: string,
}

export function FailoverButton(props: Props) {
    const {request, cluster, disabled, name} = props

    const failover = useRouterNodeFailover(cluster)
    // NOTE: in patroni we cannot use host for leader and candidate, we need to send patroni.name
    const body = {candidate: name}

    return (
        <Access feature={Feature.ManageNodeDbFailover}>
            <AlertButton
                color={"error"}
                size={"small"}
                label={"Failover"}
                title={`Make a failover of ${request.connection.host}?`}
                description={`It will failover to current node of postgres, that will cause some downtime 
                and potential data loss. Usually it is recommended to use switchover, but if you don't have a
                leader you won't be able to do switchover and here failover can be useful.`}
                disabled={disabled}
                loading={failover.isPending}
                onClick={() => failover.mutate({...request, body})}
            />
        </Access>
    )
}
