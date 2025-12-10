import {Box, Button} from "@mui/material"
import {useState} from "react"

import {useRouterQueryCancel, useRouterQueryTerminate} from "../../../api/query/hook"
import {QueryConnection} from "../../../api/query/type"
import {SxPropsMap} from "../../../app/type"
import {MenuButton} from "../../view/button/MenuButton"

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "space-evenly", color: "text.secondary", padding: "0 3px"},
    actionButton: {padding: "2px 4px", fontSize: "10px"},
}

type Props = {
    connection: QueryConnection,
    refetch: () => void,
    pid: number,
}

export function QueryTableActions(props: Props) {
    const {connection, refetch, pid} = props
    const [open, setOpen] = useState(false)

    const terminate = useRouterQueryTerminate(refetch)
    const cancel = useRouterQueryCancel(refetch)

    return (
        <MenuButton open={open} onChange={(v) => setOpen(v)}>
            <Box sx={SX.box}>
                <Button
                    sx={SX.actionButton}
                    size={"small"}
                    variant={"text"}
                    color={"error"}
                    onClick={() => {terminate.mutate({connection, pid}); setOpen(false)}}
                >
                    Terminate
                </Button>
                <Button
                    sx={SX.actionButton}
                    size={"small"}
                    variant={"text"}
                    onClick={() => {cancel.mutate({connection, pid}); setOpen(false)}}
                >
                    Cancel
                </Button>
            </Box>
        </MenuButton>
    )
}
