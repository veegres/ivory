import {ErrorSshMissing} from "../../../../shared/component/box/ErrorManual"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {Feature} from "../../../Feature"
import {ManageAccessBox} from "../../../management/component/ManageAccess"
import {PlatformVaultConnection} from "../../api/NodeType"
import {PlatformInfo} from "./PlatformInfo"
import {PlatformMetrics} from "./PlatformMetrics"
import {PlatformOverview} from "./PlatformOverview"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1, backgroundImage: "inherit", backgroundColor: "inherit"},
}

type Props = {
    connection?: PlatformVaultConnection,
}

export function Platform(props: Props) {
    const {connection} = props
    if (!connection) return <ErrorSshMissing/>
    return (
        <ManageAccessBox sx={SX.box} feature={Feature.ViewNodePlatform} error={true}>
            <PlatformInfo connection={connection}/>
            <PlatformMetrics connection={connection}/>
            <PlatformOverview connection={connection}/>
        </ManageAccessBox>
    )
}