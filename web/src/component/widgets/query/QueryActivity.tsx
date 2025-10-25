import {Box} from "@mui/material";

import {useRouterActivity} from "../../../api/query/hook";
import {QueryConnection} from "../../../api/query/type";
import {SxPropsMap} from "../../../app/type";
import {RefreshIconButton} from "../../view/button/IconButtons";
import {QueryTable} from "./QueryTable";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1, width: "320px", minHeight: "250px"},
    button: {position: "absolute", right: 10, top: 0, bottom: 0, display: "flex", alignItems: "center"},
    label: {position: "relative", textAlign: "center", fontSize: "15px", fontWeight: "bold", color: "text.secondary"},
    help: {fontSize: "9px", fontWeight: "normal", color: "text.disabled"},
    info: {
        display: "flex", justifyContent: "center", alignItems: "center", padding: "0 15px", height: "45px",
        fontSize: "11px", color: "text.secondary", textAlign: "center",
    },
    error: {color: "error.light"},
}

type Props = {
    connection: QueryConnection,
}

export function QueryActivity(props: Props) {
    const {connection} = props
    const {data, isError, error, refetch, isFetching} = useRouterActivity(connection)
    const table = isError ? undefined : data
    return (
        <Box sx={SX.box}>
            <Box sx={SX.label}>
                <Box>Your Active Queries</Box>
                <Box sx={SX.help}>[ hold shift for horizontal scrolling ]</Box>
                <Box sx={SX.button}>
                    <RefreshIconButton
                        disabled={isFetching}
                        onClick={() => refetch()}
                    />
                </Box>
            </Box>
            <QueryTable
                connection={connection}
                queryKey={"activity"}
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
