import {Box, Button} from "@mui/material"
import {useState} from "react"

import {MenuButton} from "../../../shared/component/button/MenuButton"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {Feature} from "../../Feature"
import {ManageAccess} from "../../management/component/ManageAccess"
import {useRouterQueryCancel, useRouterQueryTerminate} from "../api/QueryHook"
import {Connection} from "../api/QueryType"

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "space-evenly", alignItems: "center", color: "text.secondary", gap: "1px", height: "22px"},
    actionButton: {height: "22px", padding: "2px 6px 0px", fontSize: "10px", lineHeight: 1},
}

type Props = {
    connection: Connection,
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
                <ManageAccess feature={Feature.ManageQueryDbTerminate}>
                    <Button
                        sx={SX.actionButton}
                        size={"small"}
                        variant={"text"}
                        color={"error"}
                        onClick={() => {terminate.mutate({connection, pid}); setOpen(false)}}
                    >
                        Terminate
                    </Button>
                </ManageAccess>
                <ManageAccess feature={Feature.ManageQueryDbCancel}>
                    <Button
                        sx={SX.actionButton}
                        size={"small"}
                        variant={"text"}
                        onClick={() => {cancel.mutate({connection, pid}); setOpen(false)}}
                    >
                        Cancel
                    </Button>
                </ManageAccess>
            </Box>
        </MenuButton>
    )
}
