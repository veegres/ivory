import {ErrorSshMissing} from "../../../../shared/component/box/ErrorManual"
import {Feature} from "../../../Feature"
import {ManageAccess} from "../../../management/component/ManageAccess"
import {PlatformVaultConnection} from "../../api/NodeType"
import {PlatformInfo} from "./PlatformInfo"
import {PlatformMetrics} from "./PlatformMetrics"
import {PlatformOverview} from "./PlatformOverview"

type Props = {
    connection?: PlatformVaultConnection,
}

export function Platform(props: Props) {
    const {connection} = props
    if (!connection) return <ErrorSshMissing/>
    return (
        <ManageAccess feature={Feature.ViewNodePlatform} error={true}>
            <PlatformInfo connection={connection}/>
            <PlatformMetrics connection={connection}/>
            <PlatformOverview connection={connection}/>
        </ManageAccess>
    )
}