import {Box} from "@mui/material"

import {ConnectionRequest} from "../../../api/postgres"
import {useRouterActivity} from "../../../api/query/hook"
import {QueryApi} from "../../../api/query/router"
import {SxPropsMap} from "../../../app/type"
import {Refresher} from "../refresher/Refresher"
import {QueryTable} from "./QueryTable"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", width: "320px", minHeight: "250px"},
    head: {display: "flex", justifyContent: "space-between", alignItems: "center", padding: "0px 10px"},
    label: {fontSize: "15px", fontWeight: "bold", color: "text.secondary", textTransform: "uppercase", textAlign: "center"},
    help: {fontSize: "9px", fontWeight: "normal", color: "text.disabled", textAlign: "center"},
    info: {
        display: "flex", justifyContent: "center", alignItems: "center", padding: "0 15px", height: "45px",
        fontSize: "11px", color: "text.secondary", textAlign: "center",
    },
    error: {color: "error.light"},
}

type Props = {
    connection: ConnectionRequest,
}

export function QueryActivity(props: Props) {
    const {connection} = props
    const {data, isError, error, refetch} = useRouterActivity(connection)
    const table = isError ? undefined : data
    return (
        <Box sx={SX.box}>
            <Box>
                <Box sx={SX.head}>
                    <Box sx={SX.label}>Active Queries</Box>
                    <Box><Refresher queryKeys={[QueryApi.activity.key()]} size={26} defaultPeriod={["5s", 5000]}/></Box>
                </Box>
                <Box sx={SX.help}>[ hold shift for horizontal scrolling ]</Box>
            </Box>
            <QueryTable
                connection={connection}
                refetch={refetch}
                height={200}
                width={320}
                data={table}
                showIndexColumn={false}
            />
            <Box sx={SX.info}>
                {isError ? (
                    <Box sx={SX.error}>{error.message}</Box>
                ) : (
                    <Box>This section displays the queries submitted within your current session in Ivory</Box>
                )}
            </Box>
        </Box>
    )
}
