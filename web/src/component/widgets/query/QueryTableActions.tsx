import {Box, Button} from "@mui/material"
import {useState} from "react"

import {Permission} from "../../../api/permission/type"
import {ConnectionRequest} from "../../../api/postgres"
import {useRouterQueryCancel, useRouterQueryTerminate} from "../../../api/query/hook"
import {SxPropsMap} from "../../../app/type"
import {MenuButton} from "../../view/button/MenuButton"
import {Access} from "../access/Access"

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "space-evenly", alignItems: "center", color: "text.secondary", padding: "0 3px", height: "22px"},
    actionButton: {padding: "2px 4px", fontSize: "10px"},
}

type Props = {
    connection: ConnectionRequest,
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
                <Access permission={Permission.ManageQueryExecuteTerminate}>
                    <Button
                        sx={SX.actionButton}
                        size={"small"}
                        variant={"text"}
                        color={"error"}
                        onClick={() => {terminate.mutate({connection, pid}); setOpen(false)}}
                    >
                        Terminate
                    </Button>
                </Access>
                <Access permission={Permission.ManageQueryExecuteCancel}>
                    <Button
                        sx={SX.actionButton}
                        size={"small"}
                        variant={"text"}
                        onClick={() => {cancel.mutate({connection, pid}); setOpen(false)}}
                    >
                        Cancel
                    </Button>
                </Access>
            </Box>
        </MenuButton>
    )
}
