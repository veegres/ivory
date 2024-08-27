import {QueryConnection} from "../../../type/query";
import {useRouterActivity} from "../../../router/query";
import {QueryTable} from "./QueryTable";
import {Box} from "@mui/material";
import {SxPropsMap} from "../../../type/general";

const SX: SxPropsMap = {
    label: {textAlign: "center", padding: "5px 10px", fontSize: "15px", fontWeight: "bold", color: "text.secondary"},
    help: {fontSize: "9px", fontWeight: "normal", color: "text.disabled"},
}

type Props = {
    connection: QueryConnection,
}

export function QueryActivity(props: Props) {
    const {connection} = props
    const {data, error} = useRouterActivity(connection)
    return (
        <Box>
            <Box sx={SX.label}>
                <Box>User Active Queries</Box>
                <Box sx={SX.help}>[ hold shift for horizontal scrolling ]</Box>
            </Box>
            <QueryTable
                connection={connection}
                height={200}
                width={320}
                data={data}
                error={error}
                showIndexColumn={false}
            />
        </Box>
    )
}
