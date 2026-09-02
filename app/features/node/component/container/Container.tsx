import {ErrorSshMissing} from "../../../../shared/component/box/ErrorManual"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {Feature} from "../../../Feature"
import {ManageAccessBox} from "../../../management/component/ManageAccess"
import {PlatformVaultConnection} from "../../api/NodeType"
import {ContainerHead} from "./ContainerHead"
import {ContainerMetrics} from "./ContainerMetrics"
import {ContainerOverview} from "./ContainerOverview"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
}

type Props = {
    connection?: PlatformVaultConnection,
    name: string,
}

export function Container(props: Props) {
    const {connection, name} = props
    if (!connection) return <ErrorSshMissing/>
    return (
        <ManageAccessBox sx={SX.box} feature={Feature.ViewNodePlatformContainer} error={true}>
            <ContainerHead connection={connection} name={name}/>
            <ContainerMetrics connection={connection} name={name}/>
            <ContainerOverview connection={connection} name={name}/>
        </ManageAccessBox>
    )
}
