import {ErrorSshMissing} from "../../../../shared/component/box/ErrorManual"
import {Feature} from "../../../feature"
import {ManageAccess} from "../../../management/component/ManageAccess"
import {PlatformConnection} from "../../api/type"
import {PlatformMetrics} from "./PlatformMetrics"

type Props = {
    connection?: PlatformConnection,
}

export function Platform(props: Props) {
    const {connection} = props
    if (!connection) return <ErrorSshMissing/>
    return (
        <ManageAccess feature={Feature.ViewNodePlatform} error={true}>
            <PlatformMetrics connection={connection}/>
        </ManageAccess>
    )
}