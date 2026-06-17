import {Box} from "@mui/material"

import {ErrorSshMissing} from "../../../../shared/component/box/ErrorManual"
import {SxPropsMap} from "../../../../shared/helper/type"
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
        <Box sx={SX.box}>
            <ContainerOverview connection={connection}/>
        </Box>
    )
}
