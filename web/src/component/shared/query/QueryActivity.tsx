import {QueryConnection} from "../../../type/query";
import {useRouterActivity} from "../../../router/query";
import {QueryTable} from "./QueryTable";
import {Box} from "@mui/material";
import {SxPropsMap} from "../../../type/general";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1, width: "320px", minHeight: "250px"},
    label: {textAlign: "center", fontSize: "15px", fontWeight: "bold", color: "text.secondary"},
    help: {fontSize: "9px", fontWeight: "normal", color: "text.disabled"},
}

type Props = {
    connection: QueryConnection,
}

export function QueryActivity(props: Props) {
    const {connection} = props
    const {data} = useRouterActivity(connection)
    return (
        <Box sx={SX.box}>
            <Box sx={SX.label}>
                <Box>User Active Queries</Box>
                <Box sx={SX.help}>[ hold shift for horizontal scrolling ]</Box>
            </Box>
            <QueryTable
                connection={connection}
                queryKey={"activity"}
                height={200}
                width={320}
                data={data}
                showIndexColumn={false}
            />
        </Box>
    )
}
