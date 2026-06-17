import {AlertButton} from "../../../../shared/component/button/AlertButton"
import {Feature} from "../../../feature"
import {ManageAccess} from "../../../management/component/ManageAccess"
import {useRouterNodeReload} from "../../api/hook"
import {KeeperOneRequest} from "../../api/type"

type Props = {
    request: KeeperOneRequest,
    cluster: string,
}

export function KeeperReloadButton(props: Props) {
    const {request, cluster} = props
    const reload = useRouterNodeReload(cluster)

    return (
        <ManageAccess feature={Feature.ManageNodeDbReload}>
            <AlertButton
                size={"small"}
                label={"Reload"}
                title={`Make a reload of ${request.host}?`}
                description={`It will reload postgres config, it doesn't have any downtime. It won't help if pending 
                restart is true, some parameters require postgres restart.`}
                loading={reload.isPending}
                onClick={() => reload.mutate(request)}
            />
        </ManageAccess>
    )
}
