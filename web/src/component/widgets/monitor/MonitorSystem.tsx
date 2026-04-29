import {Box} from "@mui/material"

import {useRouterNodeMetrics} from "../../../api/node/hook"
import {Connection as NodeConnection} from "../../../api/node/type"
import {ErrorSmart} from "../../view/box/ErrorSmart"

type Props = {
    connection: NodeConnection,
}

export function MonitorSystem(props: Props) {
    const metrics = useRouterNodeMetrics({connection: props.connection}, true)

    if (metrics.error) return <ErrorSmart error={metrics.error}/>

    return (
        <Box>{JSON.stringify(metrics.data)}</Box>
    )
}