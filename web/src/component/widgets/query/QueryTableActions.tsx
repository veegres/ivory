import {Box, Button} from "@mui/material";
import {SxPropsMap} from "../../../api/management/type";
import {useRouterQueryCancel, useRouterQueryTerminate} from "../../../api/query/hook";
import {QueryConnection} from "../../../api/query/type";

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "space-evenly", color: "text.secondary", padding: "0 3px"},
    actionButton: {padding: "2px 4px", fontSize: "10px"},
}

type Props = {
    connection: QueryConnection,
    queryKey: string,
    pid: number,
}

export function QueryTableActions(props: Props) {
    const {connection, queryKey, pid} = props
    const terminate = useRouterQueryTerminate(queryKey)
    const cancel = useRouterQueryCancel(queryKey)

    return (
        <Box sx={SX.box}>
            <Button
                sx={SX.actionButton}
                size={"small"}
                variant={"text"}
                color={"error"}
                onClick={() => terminate.mutate({connection, pid})}
            >
                Terminate
            </Button>
            <Button
                sx={SX.actionButton}
                size={"small"}
                variant={"text"}
                onClick={() => cancel.mutate({connection, pid})}
            >
                Cancel
            </Button>
        </Box>
    )
}
