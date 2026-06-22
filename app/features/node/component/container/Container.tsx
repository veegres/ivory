import {ErrorSshMissing} from "../../../../shared/component/box/ErrorManual"
import {SxPropsMap} from "../../../../shared/helper/type"
import {Feature} from "../../../feature"
import {ManageAccessBox} from "../../../management/component/ManageAccess"
import {PlatformConnection} from "../../api/type"
import {ContainerOverview} from "./ContainerOverview"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", justifyContent: "center", gap: 2},
}

type Props = {
    connection?: PlatformConnection,
}

export function Container(props: Props) {
    const {connection} = props
    if (!connection) return <ErrorSshMissing/>
    return (
        <ManageAccessBox sx={SX.box} feature={Feature.ViewNodePlatformContainer} displayError={true}>
            <ContainerOverview connection={connection}/>
        </ManageAccessBox>
    )
}
