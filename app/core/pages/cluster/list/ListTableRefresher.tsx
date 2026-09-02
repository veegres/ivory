import {ClusterApi} from "../../../../features/cluster/api/ClusterRouter"
import {Refresher} from "../../../widgets/browser/Refresher"

type Props = {
    size?: number,
    width?: number,
}

export function ListTableRefresher(props: Props) {
    return (
        <Refresher
            size={props.size}
            width={props.width}
            queryKeys={[ClusterApi.list.keyCommon(), ClusterApi.overview.keyCommon()]}
        />
    )
}
