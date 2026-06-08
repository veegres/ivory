import {Feature} from "../../../features/feature"
import {useRouterNodeReload} from "../../../features/node/hook"
import {KeeperOneRequest} from "../../../features/node/type"
import {AlertButton} from "../../../shared/component/button/AlertButton"
import {Access} from "../access/Access"

type Props = {
    request: KeeperOneRequest,
    cluster: string,
}

export function ReloadButton(props: Props) {
    const {request, cluster} = props
    const reload = useRouterNodeReload(cluster)

    return (
        <Access feature={Feature.ManageNodeDbReload}>
            <AlertButton
                size={"small"}
                label={"Reload"}
                title={`Make a reload of ${request.host}?`}
                description={`It will reload postgres config, it doesn't have any downtime. It won't help if pending 
                restart is true, some parameters require postgres restart.`}
                loading={reload.isPending}
                onClick={() => reload.mutate(request)}
            />
        </Access>
    )
}
