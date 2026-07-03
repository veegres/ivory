import {ErrorSshMissing} from "../../../../shared/component/box/ErrorManual"
import {Feature} from "../../../feature"
import {ManageAccess} from "../../../management/component/ManageAccess"
import {PlatformVaultConnection} from "../../api/type"
import {PlatformLogs} from "./PlatformLogs"
import {PlatformMetrics} from "./PlatformMetrics"

type Props = {
    connection?: PlatformVaultConnection,
}

export function Platform(props: Props) {
    const {connection} = props
    if (!connection) return <ErrorSshMissing/>
    return (
        <ManageAccess feature={Feature.ViewNodePlatform} error={true}>
            <PlatformMetrics connection={connection}/>
            <PlatformLogs connection={connection}/>
        </ManageAccess>
    )
}