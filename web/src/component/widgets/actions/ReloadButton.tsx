import {Feature} from "../../../api/feature"
import {useRouterNodeReload} from "../../../api/node/hook"
import {KeeperConnection} from "../../../api/node/type"
import {AlertButton} from "../../view/button/AlertButton"
import {Access} from "../access/Access"

type Props = {
    con: KeeperConnection,
    cluster: string,
}

export function ReloadButton(props: Props) {
    const {con, cluster} = props
    const reload = useRouterNodeReload(cluster)

    return (
        <Access feature={Feature.ManageNodeDbReload}>
            <AlertButton
                size={"small"}
                label={"Reload"}
                title={`Make a reload of ${con.host}?`}
                description={`It will reload postgres config, it doesn't have any downtime. It won't help if pending 
                restart is true, some parameters require postgres restart.`}
                loading={reload.isPending}
                onClick={() => reload.mutate(con)}
            />
        </Access>
    )
}
