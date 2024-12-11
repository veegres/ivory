import {QueryConnection} from "../../../type/query";
import {useRouterActivity} from "../../../router/query";
import {QueryTable} from "./QueryTable";
import {Box} from "@mui/material";
import {SxPropsMap} from "../../../type/general";
import {RefreshIconButton} from "../../view/button/IconButtons";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1, width: "320px", minHeight: "250px"},
    button: {position: "absolute", right: 10, top: 0, bottom: 0, display: "flex", alignItems: "center"},
    label: {position: "relative", textAlign: "center", fontSize: "15px", fontWeight: "bold", color: "text.secondary"},
    help: {fontSize: "9px", fontWeight: "normal", color: "text.disabled"},
}

type Props = {
    connection: QueryConnection,
}

export function QueryActivity(props: Props) {
    const {connection} = props
    const {data, refetch, isFetching} = useRouterActivity(connection)
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
                data={data}
                showIndexColumn={false}
            />
        </Box>
    )
}
