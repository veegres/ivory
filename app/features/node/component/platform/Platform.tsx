import {ErrorSshMissing} from "../../../../shared/component/box/ErrorManual"
import {PlatformConnection} from "../../api/type"
import {PlatformMetrics} from "./PlatformMetrics"

type Props = {
    connection?: PlatformConnection,
}

export function Platform(props: Props) {
    const {connection} = props
    if (!connection) return <ErrorSshMissing/>
    return (
        <PlatformMetrics connection={connection}/>
    )
}