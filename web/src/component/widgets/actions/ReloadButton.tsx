import {useRouterNodeReload} from "../../../api/node/hook"
import {NodeRequest} from "../../../api/node/type"
import {Permission} from "../../../api/permission/type"
import {AlertButton} from "../../view/button/AlertButton"
import {Access} from "../access/Access"

type Props = {
    request: NodeRequest,
    cluster: string,
}

export function ReloadButton(props: Props) {
    const {request, cluster} = props
    const reload = useRouterNodeReload(cluster)

    return (
        <Access permission={Permission.ManageNodeReload}>
            <AlertButton
                size={"small"}
                label={"Reload"}
                title={`Make a reload of ${request.sidecar.host}?`}
                description={`It will reload postgres config, it doesn't have any downtime. It won't help if pending 
                restart is true, some parameters require postgres restart.`}
                loading={reload.isPending}
                onClick={() => reload.mutate(request)}
            />
        </Access>
    )
}
