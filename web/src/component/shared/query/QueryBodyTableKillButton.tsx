import {Box, Button} from "@mui/material";
import {SxPropsMap} from "../../../type/general";
import {useRouterQueryCancel, useRouterQueryTerminate} from "../../../router/query";
import {QueryConnection} from "../../../type/query";

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "space-evenly", color: "text.secondary", padding: "0 3px"},
    actionButton: {padding: "0px 4px", fontSize: "10px"},
}

type Props = {
    connection: QueryConnection,
    queryUuid?: string,
    pid: number,
}

export function QueryBodyTableKillButton(props: Props) {
    const {connection, queryUuid, pid} = props
    // TODO fix rerender we not always have Uuid, check (console and activity)
    const terminate = useRouterQueryTerminate(queryUuid)
    const cancel = useRouterQueryCancel(queryUuid)

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
