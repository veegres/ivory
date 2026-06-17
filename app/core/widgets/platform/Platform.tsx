import {PlatformConnection} from "../../../features/node/type"
import {ErrorSshMissing} from "../../../shared/component/box/ErrorManual"
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