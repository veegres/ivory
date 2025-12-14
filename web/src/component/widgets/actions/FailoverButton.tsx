import {useRouterInstanceFailover} from "../../../api/instance/hook"
import {InstanceRequest} from "../../../api/instance/type"
import {Permission} from "../../../api/permission/type"
import {AlertButton} from "../../view/button/AlertButton"
import {Access} from "../access/Access"

type Props = {
    request: InstanceRequest,
    cluster: string,
    disabled: boolean,
}

export function FailoverButton(props: Props) {
    const {request, cluster, disabled} = props

    const failover = useRouterInstanceFailover(cluster)
    // NOTE: in patroni we cannot use host for leader and candidate, we need to send patroni.name
    const body = {candidate: request.sidecar.name}

    return (
        <Access permission={Permission.ManageInstanceFailover}>
            <AlertButton
                color={"error"}
                size={"small"}
                label={"Failover"}
                title={`Make a failover of ${request.sidecar.host}?`}
                description={`It will failover to current instance of postgres, that will cause some downtime 
                and potential data loss. Usually it is recommended to use switchover, but if you don't have a
                leader you won't be able to do switchover and here failover can be useful.`}
                disabled={disabled}
                loading={failover.isPending}
                onClick={() => failover.mutate({...request, body})}
            />
        </Access>
    )
}
