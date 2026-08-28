import {ErrorSshMissing} from "../../../../shared/component/box/ErrorManual"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {Feature} from "../../../Feature"
import {ManageAccessBox} from "../../../management/component/ManageAccess"
import {PlatformVaultConnection} from "../../api/NodeType"
import {SystemInfo} from "./SystemInfo"
import {SystemMetrics} from "./SystemMetrics"
import {SystemOverview} from "./SystemOverview"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1, backgroundImage: "inherit", backgroundColor: "inherit"},
}

type Props = {
    connection?: PlatformVaultConnection,
}

export function System(props: Props) {
    const {connection} = props
    if (!connection) return <ErrorSshMissing/>
    return (
        <ManageAccessBox sx={SX.box} feature={Feature.ViewNodeSystem} error={true}>
            <SystemInfo connection={connection}/>
            <SystemMetrics connection={connection}/>
            <SystemOverview connection={connection}/>
        </ManageAccessBox>
    )
}