import {Box, Button} from "@mui/material"
import {cloneElement} from "react"

import {SxPropsMap} from "../../../shared/helper/HelperType"
import {PermissionOptions} from "../../../shared/helper/HelperUtils"
import {Feature} from "../../Feature"
import {Status} from "../api/PermissionType"
import {PermissionButtons} from "./PermissionButtons"

const SX: SxPropsMap = {
    item: {
        display: "flex", justifyContent: "space-between", alignItems: "center",
        borderTop: 1, borderColor: "divider", padding: "4px 8px", height: "35px",
    },
    button: {padding: "2px 5px"},
    wrap: {display: "flex", alignItems: "center", gap: 1},
}

type Props = {
    name: Feature,
    status: Status,
    username: string,
    view?: "admin" | "user",
}

export function PermissionListItem(props: Props) {
    const {username, name, status, view = "user"} = props
    const options = PermissionOptions[status]

    return (
        <Box sx={SX.item}>
            <Box sx={SX.wrap}>
                {cloneElement(options.icon, {sx: {color: options.color, fontSize: "20px"}})}
                <Box>{name}</Box>
            </Box>
            {renderButton()}
        </Box>
    )

    function renderButton() {
        if (view === "admin") {
            return <PermissionButtons
                username={username}
                permissions={[[name, status]]}
                approve={status === Status.NOT_PERMITTED || status === Status.PENDING}
                reject={status === Status.GRANTED || status === Status.PENDING}
            />
        }
        if (status === Status.NOT_PERMITTED) {
            return <PermissionButtons username={username} permissions={[[name, status]]} request={true}/>
        }
        return renderStatus(status)
    }

    function renderStatus(status: Status) {
        return <Button sx={SX.button} disabled={true}>{Status[status]}</Button>
    }


}