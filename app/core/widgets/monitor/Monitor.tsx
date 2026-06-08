import {Box} from "@mui/material"

import {PlatformConnection} from "../../../features/node/type"
import {ErrorSshMissing} from "../../../shared/component/box/ErrorManual"
import {SxPropsMap} from "../../../shared/helper/type"
import {MonitorContainer} from "./MonitorContainer"
import {MonitorPlatform} from "./MonitorPlatform"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", justifyContent: "center", gap: 2},
}

type Props = {
    connection?: PlatformConnection,
}

export function Monitor(props: Props) {
    const {connection} = props

    if (!connection) return <ErrorSshMissing/>

    return (
        <Box sx={SX.box}>
            <MonitorPlatform connection={connection}/>
            <MonitorContainer connection={connection}/>
        </Box>
    )
}
